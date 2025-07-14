package handler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2" 

	"IoTT/internal/database"
	"IoTT/internal/models"
	"IoTT/internal/telegram"
	"IoTT/internal/worker"

)



func HandleSensorData(c *fiber.Ctx) error {
	
	var payloads []models.SensorPayload
	if err := c.BodyParser(&payloads); err != nil {
		
		var singlePayload models.SensorPayload
		if errSingle := c.BodyParser(&singlePayload); errSingle != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON payload. Expected an object or an array of objects.",
			})
		}
		payloads = append(payloads, singlePayload)
	}

	if len(payloads) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Received empty data list.",
		})
	}

	db := database.GetDB()
	if db == nil {
		log.Println("Error: Database connection is not initialized.")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection not initialized",
		})
	}

	
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting database transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
	}
	
	defer tx.Rollback()

	var tempsToBatch []worker.TempBatchData
	var rhsToBatch []worker.RhBatchData
	var proxsToBatch []worker.ProxBatchData
	processedItemCount := 0

	for _, p := range payloads {
		sensorNoVal := 1 
		if p.No != nil {
			sensorNoVal = *p.No
		}
		areaIDVal := 0
		if p.AreaID != nil {
			areaIDVal = *p.AreaID
		}
		doorIDVal := 0
		if p.DoorID != nil {
			doorIDVal = *p.DoorID
		}

		var parsedTS time.Time
		var tsVal string
		if p.TS != nil && *p.TS != "" {
			tsVal = *p.TS
			parsedTS, _ = time.Parse(time.RFC3339Nano, tsVal) 
		} else {
			parsedTS = time.Now().UTC()
			tsVal = parsedTS.Format(time.RFC3339Nano)
		}

		if p.Temp != nil {
			tempsToBatch = append(tempsToBatch, worker.TempBatchData{Value: *p.Temp, AreaID: areaIDVal, No: sensorNoVal, TS: tsVal})
			sensorKey := fmt.Sprintf("temp-%d-%d", areaIDVal, sensorNoVal)
			models.RegisterOrUpdateSensorStatus(sensorKey, "temp", areaIDVal, sensorNoVal, 0, parsedTS)
		}

		if p.RH != nil {
			rhsToBatch = append(rhsToBatch, worker.RhBatchData{Value: *p.RH, AreaID: areaIDVal, No: sensorNoVal, TS: tsVal})
			sensorKey := fmt.Sprintf("rh-%d-%d", areaIDVal, sensorNoVal)
			models.RegisterOrUpdateSensorStatus(sensorKey, "rh", areaIDVal, sensorNoVal, 0, parsedTS)
		}

		if p.Prox != nil {
			proxsToBatch = append(proxsToBatch, worker.ProxBatchData{Value: *p.Prox, DoorID: doorIDVal, TS: tsVal})
			sensorKey := fmt.Sprintf("prox-%d", doorIDVal)
			models.RegisterOrUpdateSensorStatus(sensorKey, "prox", 0, 0, doorIDVal, parsedTS)
		}
		processedItemCount++
	}

	var wg sync.WaitGroup
	
	errs := make(chan error, 3)

	if len(tempsToBatch) > 0 {
		wg.Add(1)
		go func(data []worker.TempBatchData) {
			defer wg.Done()
			
			if err := worker.BatchInsertTemp(tx, data); err != nil {
				errs <- fmt.Errorf("gagal batch insert temp: %w", err)
				return 
			}
			
			for _, d := range data {
				safetyStatus := worker.EvaluateTemp(d.AreaID, d.No, d.Value)
				if safetyStatus.IsAlert && safetyStatus.Message != "" {
					telegram.SendAlert(fmt.Sprintf("🔥 **Alert Temperature**\n📍 Area %d Titik %d  %s", d.AreaID, d.No, safetyStatus.Message))
				}
			}
		}(tempsToBatch)
	}

	if len(rhsToBatch) > 0 {
		wg.Add(1)
		go func(data []worker.RhBatchData) {
			defer wg.Done()
			
			if err := worker.BatchInsertRh(tx, data); err != nil {
				errs <- fmt.Errorf("gagal batch insert rh: %w", err)
				return 
			}
			
			for _, d := range data {
				safetyStatus := worker.EvaluateRh(d.AreaID, d.No, d.Value)
				if safetyStatus.IsAlert && safetyStatus.Message != "" {
					telegram.SendAlert(fmt.Sprintf("💧 **Alert Humidity**\n📍Area %d Titik %d %s", d.AreaID, d.No, safetyStatus.Message))
				}
			}
		}(rhsToBatch)
	}

	if len(proxsToBatch) > 0 {
		wg.Add(1)
		go func(data []worker.ProxBatchData) {
			defer wg.Done()
			
			if err := worker.BatchInsertProx(tx, data); err != nil {
				errs <- fmt.Errorf("gagal batch insert prox: %w", err)
			}
		}(proxsToBatch)
	}

	wg.Wait()
	close(errs) 

	
	var errorMessages []string
	for err := range errs {
		log.Printf("Error selama pemrosesan data sensor: %v", err)
		errorMessages = append(errorMessages, err.Error())
	}

	if len(errorMessages) > 0 {
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Terjadi kesalahan saat menyimpan sebagian atau seluruh data.",
			"errors":  errorMessages,
		})
	}

	
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal menyimpan data ke database setelah pemrosesan.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": fmt.Sprintf("Successfully processed %d data items.", processedItemCount),
	})
}


func HandleTelegramWebhook(c *fiber.Ctx) error {
	var update interface{}
	
	if err := c.BodyParser(&update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid update format",
		})
	}

	log.Printf("Received Telegram update: %+v", update)

	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}
