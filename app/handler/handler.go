package handler

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xuri/excelize/v2"

	"IoTT/app/models"
	"IoTT/internal/config"
	"IoTT/internal/database"
	internalmodels "IoTT/internal/models"
	"IoTT/internal/worker"
)

const getTempDataByAreaAndRangeQuery = `
    SELECT no, value, ts 
    FROM temp 
    WHERE area_id = ? AND ts >= ? AND ts < ? 
    ORDER BY ts ASC;`

const getRhDataByAreaAndRangeQuery = `
    SELECT no, value, ts 
    FROM rh 
    WHERE area_id = ? AND ts >= ? AND ts < ? 
    ORDER BY ts ASC;`

const getTempDataBySensorAndRangeQuery = `
    SELECT no, value, ts 
    FROM temp 
    WHERE area_id = ? AND no = ? AND ts >= ? AND ts < ? 
    ORDER BY ts ASC;`

const getRhDataBySensorAndRangeQuery = `
    SELECT no, value, ts 
    FROM rh 
    WHERE area_id = ? AND no = ? AND ts >= ? AND ts < ? 
    ORDER BY ts ASC;`

const getAreaSummaryStatsBySensorQuery = `
-- Versi query yang diperbaiki dan stabil, placeholder diubah ke '?'
SELECT
    (SELECT AVG(value) FROM temp WHERE area_id = ? AND no = ? AND ts >= ? AND ts < ?) as avg_temp,
    (SELECT AVG(value) FROM rh   WHERE area_id = ? AND no = ? AND ts >= ? AND ts < ?) as avg_rh,
    (SELECT MAX(value) FROM temp WHERE area_id = ? AND no = ? AND ts >= ? AND ts < ?) as max_temp,
    (SELECT MIN(value) FROM temp WHERE area_id = ? AND no = ? AND ts >= ? AND ts < ?) as min_temp,
    (SELECT MAX(value) FROM rh   WHERE area_id = ? AND no = ? AND ts >= ? AND ts < ?) as max_rh,
    (SELECT MIN(value) FROM rh   WHERE area_id = ? AND no = ? AND ts >= ? AND ts < ?) as min_rh,
    (SELECT value FROM temp WHERE area_id = ? AND no = ? AND ts >= ? AND ts < ? ORDER BY ts DESC LIMIT 1) as last_temp,
    (SELECT value FROM rh   WHERE area_id = ? AND no = ? AND ts >= ? AND ts < ? ORDER BY ts DESC LIMIT 1) as last_rh;
`

const getLatestAndPreviousSensorDataQuery = `
-- Query diubah untuk menggantikan DISTINCT ON (PostgreSQL) dengan ROW_NUMBER() (SQLite).
WITH ranked_sensors AS (
    SELECT
        area_id, no, value, ts, 'temp' AS sensor_type,
        LAG(value, 1) OVER (PARTITION BY area_id, no ORDER BY ts) as prev_value
    FROM temp
    UNION ALL
    SELECT
        area_id, no, value, ts, 'rh' AS sensor_type,
        LAG(value, 1) OVER (PARTITION BY area_id, no ORDER BY ts) as prev_value
    FROM rh
),
latest_ranked AS (
    SELECT *,
           ROW_NUMBER() OVER(PARTITION BY area_id, no, sensor_type ORDER BY ts DESC) as rn
    FROM ranked_sensors
)
SELECT area_id, no, value, ts, sensor_type, prev_value
FROM latest_ranked
WHERE rn = 1
ORDER BY area_id, no, sensor_type;
`

const getCombinedSensorStatusesQuery = `
-- Query diubah untuk menggantikan FULL OUTER JOIN (PostgreSQL) dengan LEFT JOIN (SQLite).
WITH latest_temp AS (
    SELECT
        area_id, no, value, ts,
        ROW_NUMBER() OVER(PARTITION BY area_id, no ORDER BY ts DESC) as rn
    FROM temp
),
latest_rh AS (
    SELECT
        area_id, no, value, ts,
        ROW_NUMBER() OVER(PARTITION BY area_id, no ORDER BY ts DESC) as rn
    FROM rh
),
all_sensors AS (
    SELECT DISTINCT area_id, no FROM temp
    UNION
    SELECT DISTINCT area_id, no FROM rh
)
SELECT
    s.area_id,
    a.name AS area_name,
    s.no AS sensor_no,
    t.value AS temp_value,
    t.ts AS temp_ts,
    r.value AS rh_value,
    r.ts AS rh_ts
FROM all_sensors s
LEFT JOIN (SELECT * FROM latest_temp WHERE rn = 1) t ON s.area_id = t.area_id AND s.no = t.no
LEFT JOIN (SELECT * FROM latest_rh WHERE rn = 1) r ON s.area_id = r.area_id AND s.no = r.no
LEFT JOIN area a ON s.area_id = a.area_id
WHERE a.area_id IS NOT NULL
ORDER BY
    s.area_id, s.no;
`

const getCombinedSensorDataForExportQuery = `
    SELECT t.ts, t.value, r.value
    FROM temp t
    LEFT JOIN rh r ON t.area_id = r.area_id AND t.no = r.no AND t.ts = r.ts
    WHERE t.area_id = ? AND t.no = ? AND t.ts >= ? AND t.ts < ?
    ORDER BY t.ts ASC;
`

func parseAreaID(param string) (int, error) {
	if !strings.HasPrefix(param, "s0_") {
		return 0, fmt.Errorf("format ID area tidak valid, harus diawali dengan 's0_'")
	}
	idPart := strings.TrimPrefix(param, "s0_")
	areaID, err := strconv.Atoi(idPart)
	if err != nil {
		return 0, fmt.Errorf("bagian ID '%s' setelah 's0_' bukan angka yang valid", idPart)
	}
	return areaID, nil
}

func parseAreaAndDateRangeParams(c *fiber.Ctx) (areaID int, startDate, endDate time.Time, err error) {
	areaIDStr := c.Params("area_id")
	areaID, err = parseAreaID(areaIDStr)
	if err != nil {
		return
	}
	startStr := c.Query("start")
	endStr := c.Query("end")
	if startStr == "" || endStr == "" {
		err = fmt.Errorf("Query parameter 'start' dan 'end' dibutuhkan")
		return
	}
	startDate, err = time.Parse("2006-01-02", startStr)
	if err != nil {
		err = fmt.Errorf("Format tanggal 'start' tidak valid. Gunakan YYYY-MM-DD.")
		return
	}
	endDate, err = time.Parse("2006-01-02", endStr)
	if err != nil {
		err = fmt.Errorf("Format tanggal 'end' tidak valid. Gunakan YYYY-MM-DD.")
		return
	}

	endDate = endDate.AddDate(0, 0, 1)
	return
}

func ExportSensorDataToExcel(c *fiber.Ctx) error {
	areaID, startDate, endDate, err := parseAreaAndDateRangeParams(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	db := database.GetDB()
	if db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Koneksi database tidak terinisialisasi"})
	}
	fetchData := func(query string) ([]models.SensorDataPoint, error) {
		rows, err := db.Query(query, areaID, startDate, endDate)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var dataPoints []models.SensorDataPoint
		for rows.Next() {
			var dp models.SensorDataPoint
			if err := rows.Scan(&dp.No, &dp.Value, &dp.TS); err != nil {
				return nil, err
			}
			dataPoints = append(dataPoints, dp)
		}
		return dataPoints, rows.Err()
	}
	tempData, err := fetchData(getTempDataByAreaAndRangeQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data temperatur untuk ekspor"})
	}
	rhData, err := fetchData(getRhDataByAreaAndRangeQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data RH untuk ekspor"})
	}
	f := excelize.NewFile()
	defer f.Close()
	writeSheet := func(sheetName string, data []models.SensorDataPoint) {
		f.NewSheet(sheetName)
		f.SetCellValue(sheetName, "A1", "No Sensor")
		f.SetCellValue(sheetName, "B1", "Value")
		f.SetCellValue(sheetName, "C1", "Timestamp")
		for i, dp := range data {
			row := i + 2
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), dp.No)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), dp.Value)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), dp.TS.Format(time.RFC3339))
		}
	}
	writeSheet("Temperature", tempData)
	writeSheet("RH", rhData)
	f.DeleteSheet("Sheet1")
	var buffer bytes.Buffer
	if err := f.Write(&buffer); err != nil {
		log.Printf("Error writing excel file to buffer: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat file Excel"})
	}
	fileName := fmt.Sprintf("export_area_%d_%s_to_%s.xlsx", areaID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename="+fileName)
	return c.Send(buffer.Bytes())
}

func ExportSingleSensorDataToExcel(c *fiber.Ctx) error {

	areaIDStr := c.Query("area_id")
	sensorNoStr := c.Query("sensor_no")
	startStr := c.Query("start_date")
	endStr := c.Query("end_date")

	if areaIDStr == "" || sensorNoStr == "" || startStr == "" || endStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Parameter tidak lengkap (area_id, sensor_no, start_date, end_date wajib ada)"})
	}

	areaID, err := parseAreaID(areaIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	sensorNo, err := strconv.Atoi(sensorNoStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Parameter 'sensor_no' harus angka"})
	}

	startDate, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format 'start_date' salah, pakai YYYY-MM-DD"})
	}

	endDate, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format 'end_date' salah, pakai YYYY-MM-DD"})
	}

	endDate = endDate.AddDate(0, 0, 1)

	db := database.GetDB()
	if db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Koneksi database tidak terinisialisasi"})
	}

	rows, err := db.Query(getCombinedSensorDataForExportQuery, areaID, sensorNo, startDate, endDate)
	if err != nil {
		log.Printf("Error querying combined sensor data for export: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data sensor untuk ekspor"})
	}
	defer rows.Close()

	type CombinedDataPoint struct {
		TS   time.Time
		Temp sql.NullFloat64
		RH   sql.NullFloat64
	}

	var dataPoints []CombinedDataPoint
	for rows.Next() {
		var dp CombinedDataPoint
		if err := rows.Scan(&dp.TS, &dp.Temp, &dp.RH); err != nil {
			log.Printf("Error scanning combined data row for export: %v", err)
			continue
		}
		dataPoints = append(dataPoints, dp)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating combined data rows for export: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca data sensor untuk ekspor"})
	}

	f := excelize.NewFile()
	defer f.Close()
	sheetName := fmt.Sprintf("Sensor %d Area %d", sensorNo, areaID)
	f.NewSheet(sheetName)

	f.SetCellValue(sheetName, "A1", "Timestamp")
	f.SetCellValue(sheetName, "B1", "Temperature (°C)")
	f.SetCellValue(sheetName, "C1", "Humidity (%%)")

	for i, dp := range dataPoints {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), dp.TS.Format("2006-01-02 15:04:05"))
		if dp.Temp.Valid {
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), dp.Temp.Float64)
		}
		if dp.RH.Valid {
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), dp.RH.Float64)
		}
	}

	f.DeleteSheet("Sheet1")

	var buffer bytes.Buffer
	if err := f.Write(&buffer); err != nil {
		log.Printf("Error writing excel file to buffer: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat file Excel"})
	}
	fileName := fmt.Sprintf("export_sensor_%d_area_%d_%s_to_%s.xlsx", sensorNo, areaID, startDate.Format("2006-01-02"), endDate.AddDate(0, 0, -1).Format("2006-01-02"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename="+fileName)
	return c.Send(buffer.Bytes())
}

func parseSummaryRequestParams(c *fiber.Ctx) (startDate, endDate time.Time, err error) {
	startStr := c.Query("start", "today")
	endStr := c.Query("end", "today")

	now := time.Now().In(config.WIBLocation)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, config.WIBLocation)

	if strings.ToLower(startStr) == "today" {
		startDate = today
	} else {
		startDate, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			err = fmt.Errorf("format tanggal 'start' tidak valid. Gunakan YYYY-MM-DD atau 'today'")
			return
		}
	}

	if strings.ToLower(endStr) == "today" {
		endDate = today
	} else {
		endDate, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			err = fmt.Errorf("format tanggal 'end' tidak valid. Gunakan YYYY-MM-DD atau 'today'")
			return
		}
	}

	if startDate.After(endDate) {
		err = fmt.Errorf("tanggal 'start' tidak boleh setelah tanggal 'end'")
		return
	}

	endDate = endDate.AddDate(0, 0, 1)
	return
}

func getSetpointsFromConfig(thresholds []config.ThresholdConfig, areaID, sensorNo int) (minSet, maxSet *float64) {
	for _, t := range thresholds {
		if t.AreaID == areaID && t.SensorNo == sensorNo {
			if t.Config.LowerWarning != nil {
				minSet = &t.Config.LowerWarning.Limit
			} else if t.Config.LowerCritical != nil {
				minSet = &t.Config.LowerCritical.Limit
			}
			if t.Config.UpperWarning != nil {
				maxSet = &t.Config.UpperWarning.Limit
			} else if t.Config.UpperCritical != nil {
				maxSet = &t.Config.UpperCritical.Limit
			}
			return
		}
	}
	return
}

const getSensorDataForTableQuery = `
    SELECT t.ts, t.value, r.value
    FROM temp t
    LEFT JOIN rh r ON t.area_id = r.area_id AND t.no = r.no AND t.ts = r.ts -- Join berdasarkan timestamp yang sama
    WHERE t.area_id = ? AND t.no = ? AND t.ts >= ? AND t.ts < ?
    ORDER BY t.ts DESC;
`

func GetAreaSummaryBySensorHandler(c *fiber.Ctx) error {
	areaIDStr := c.Params("area_id")
	areaID, err := parseAreaID(areaIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	startDate, endDate, err := parseSummaryRequestParams(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	noStr := c.Params("no")
	sensorNo, err := strconv.Atoi(noStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nomor sensor (no) tidak valid"})
	}

	type TableRow struct {
		TitikSensor string    `json:"titik_sensor"`
		Area        string    `json:"area"`
		Timestamp   time.Time `json:"timestamp"`
		Status      string    `json:"status"`
		Temperature *float64  `json:"temperature"`
		Humidity    *float64  `json:"humidity"`
	}

	type SummaryResponse struct {
		LastTemp   *float64   `json:"last_temp"`
		MinTemp    *float64   `json:"min_temp"`
		MaxTemp    *float64   `json:"max_temp"`
		AvgTemp    *float64   `json:"avg_temp"`
		MinSetTemp *float64   `json:"min_set_temp"`
		MaxSetTemp *float64   `json:"max_set_temp"`
		LastRH     *float64   `json:"last_rh"`
		MinRH      *float64   `json:"min_rh"`
		MaxRH      *float64   `json:"max_rh"`
		AvgRH      *float64   `json:"avg_rh"`
		MinSetRH   *float64   `json:"min_set_rh"`
		MaxSetRH   *float64   `json:"max_set_rh"`
		Table      []TableRow `json:"table"`
	}

	response := SummaryResponse{
		Table: []TableRow{},
	}

	response.MinSetTemp, response.MaxSetTemp = getSetpointsFromConfig(config.TempThresholds, areaID, sensorNo)
	response.MinSetRH, response.MaxSetRH = getSetpointsFromConfig(config.RhThresholds, areaID, sensorNo)

	db := database.GetDB()
	if db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Koneksi database tidak terinisialisasi"})
	}

	var avgTemp, avgRH, maxTemp, minTemp, maxRH, minRH, lastTemp, lastRH sql.NullFloat64
	err = db.QueryRow(getAreaSummaryStatsBySensorQuery,
		areaID, sensorNo, startDate, endDate,
		areaID, sensorNo, startDate, endDate,
		areaID, sensorNo, startDate, endDate,
		areaID, sensorNo, startDate, endDate,
		areaID, sensorNo, startDate, endDate,
		areaID, sensorNo, startDate, endDate,
		areaID, sensorNo, startDate, endDate,
		areaID, sensorNo, startDate, endDate,
	).Scan(
		&avgTemp, &avgRH, &maxTemp, &minTemp,
		&maxRH, &minRH,
		&lastTemp, &lastRH)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error querying area summary stats for area %d sensor %d: %v", areaID, sensorNo, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil statistik area"})
	}

	if avgTemp.Valid {
		roundedAvg := math.Round(avgTemp.Float64*100) / 100
		response.AvgTemp = &roundedAvg
	}
	if avgRH.Valid {
		roundedAvg := math.Round(avgRH.Float64*100) / 100
		response.AvgRH = &roundedAvg
	}
	if maxTemp.Valid {
		response.MaxTemp = &maxTemp.Float64
	}
	if minTemp.Valid {
		response.MinTemp = &minTemp.Float64
	}
	if maxRH.Valid {
		response.MaxRH = &maxRH.Float64
	}
	if minRH.Valid {
		response.MinRH = &minRH.Float64
	}
	if lastTemp.Valid {
		response.LastTemp = &lastTemp.Float64
	}
	if lastRH.Valid {
		response.LastRH = &lastRH.Float64
	}

	rows, err := db.Query(getSensorDataForTableQuery, areaID, sensorNo, startDate, endDate)
	if err != nil {
		log.Printf("Error querying table data for area %d sensor %d: %v", areaID, sensorNo, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data detail untuk tabel"})
	}
	defer rows.Close()

	areaName := database.GetAreaName(areaID)
	for rows.Next() {
		var ts time.Time
		var tempVal, rhVal sql.NullFloat64
		if err := rows.Scan(&ts, &tempVal, &rhVal); err != nil {
			log.Printf("Error scanning table data row: %v", err)
			continue
		}
		row := TableRow{
			TitikSensor: fmt.Sprintf("Sensor T/H %02d", sensorNo),
			Area:        areaName,
			Timestamp:   ts,
			Status:      "running",
		}
		if tempVal.Valid {
			row.Temperature = &tempVal.Float64
		}
		if rhVal.Valid {
			row.Humidity = &rhVal.Float64
		}
		response.Table = append(response.Table, row)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetDetailedAlerts(c *fiber.Ctx) error {
	db := database.GetDB()
	if db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Koneksi database tidak terinisialisasi"})
	}
	rows, err := db.Query(getLatestAndPreviousSensorDataQuery)
	if err != nil {
		log.Printf("Error querying latest sensor data for detailed alerts: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data sensor untuk alert"})
	}
	defer rows.Close()
	var detailedAlerts []models.DetailedAlert

	now := time.Now().In(config.WIBLocation)
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, config.WIBLocation)

	for rows.Next() {
		var (
			areaID, sensorNo int
			currentValue     float64
			ts               time.Time
			sensorType       string
			previousValue    sql.NullFloat64
		)
		if err := rows.Scan(&areaID, &sensorNo, &currentValue, &ts, &sensorType, &previousValue); err != nil {
			log.Printf("Error scanning latest sensor data row for detailed alerts: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memproses data sensor untuk alert"})
		}

		if ts.Before(startOfToday) {
			continue
		}

		var status worker.SafetyStatus
		if sensorType == "temp" {
			status = worker.EvaluateTemp(areaID, sensorNo, currentValue)
		} else if sensorType == "rh" {
			status = worker.EvaluateRh(areaID, sensorNo, currentValue)
		}
		if status.IsAlert {
			alert := models.DetailedAlert{
				AreaName:    database.GetAreaName(areaID),
				SensorNo:    sensorNo,
				SensorType:  sensorType,
				ValueAfter:  currentValue,
				Description: status.Message,
				Status:      status.Severity,
				Timestamp:   ts,
			}
			if previousValue.Valid {
				alert.ValueBefore = &previousValue.Float64
			}
			detailedAlerts = append(detailedAlerts, alert)
		}
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating latest sensor data rows for detailed alerts: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error membaca data sensor untuk alert"})
	}
	if detailedAlerts == nil {
		detailedAlerts = []models.DetailedAlert{}
	}
	return c.Status(fiber.StatusOK).JSON(detailedAlerts)
}

func Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	var userID int

	if req.Username == "admin" && req.Password == "admin123" {
		userID = 1
	} else if req.Username == "su" && req.Password == "sus" {
		userID = 2
	} else {

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Username atau password salah"})
	}

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		log.Println("FATAL: JWT_SECRET_KEY tidak diatur di .env")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Konfigurasi server tidak lengkap"})
	}

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": req.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat token otentikasi"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": tokenString})
}

func Logout(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Successfully logged out"})
}

func GetChartDataHandler(c *fiber.Ctx) error {
	type ChartDataPoint struct {
		Time  time.Time `json:"time"`
		Value float64   `json:"value"`
	}
	areaIDStr := c.Query("area_id")
	sensorNoStr := c.Query("sensor_no")
	dataType := c.Query("type")
	startStr := c.Query("start_date")
	endStr := c.Query("end_date")
	if areaIDStr == "" || sensorNoStr == "" || dataType == "" || startStr == "" || endStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Parameter tidak lengkap (area_id, sensor_no, type, start_date, end_date wajib ada)"})
	}
	areaID, err := parseAreaID(areaIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	sensorNo, err := strconv.Atoi(sensorNoStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Parameter 'sensor_no' harus angka"})
	}
	startDate, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format 'start_date' salah, pakai YYYY-MM-DD"})
	}
	endDate, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format 'end_date' salah, pakai YYYY-MM-DD"})
	}

	endDate = endDate.AddDate(0, 0, 1)
	var query string
	if dataType == "temperature" {
		query = getTempDataBySensorAndRangeQuery
	} else if dataType == "humidity" {
		query = getRhDataBySensorAndRangeQuery
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Parameter 'type' hanya boleh 'temperature' atau 'humidity'"})
	}
	db := database.GetDB()
	rows, err := db.Query(query, areaID, sensorNo, startDate, endDate)
	if err != nil {
		log.Printf("Error querying chart data: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal ambil data chart"})
	}
	defer rows.Close()
	var responseData []ChartDataPoint
	for rows.Next() {
		var dp models.SensorDataPoint
		if err := rows.Scan(&dp.No, &dp.Value, &dp.TS); err != nil {
			log.Printf("Error scanning chart data row: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal proses data chart"})
		}
		point := ChartDataPoint{
			Time:  dp.TS,
			Value: dp.Value,
		}
		responseData = append(responseData, point)
	}
	if responseData == nil {
		responseData = []ChartDataPoint{}
	}
	return c.Status(fiber.StatusOK).JSON(responseData)
}

func mapSeverityToStatus(severity string) string {
	if strings.HasPrefix(severity, "KRITIS") {
		return "danger"
	}
	if strings.HasPrefix(severity, "WASPADA") {
		return "warning"
	}
	return "normal"
}

func GetSensorStatuses(c *fiber.Ctx) error {
	db := database.GetDB()
	if db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Koneksi database belum diinisialisasi"})
	}

	rows, err := db.Query(getCombinedSensorStatusesQuery)
	if err != nil {
		log.Printf("Error saat query data sensor gabungan: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data sensor"})
	}
	defer rows.Close()

	var statuses []models.CombinedSensorStatus

	now := time.Now().In(config.WIBLocation)
	offlineDuration := 10 * time.Minute

	for rows.Next() {
		var areaID, sensorNo int
		var areaName string
		var tempValue sql.NullFloat64
		var tempTS sql.NullTime
		var rhValue sql.NullFloat64
		var rhTS sql.NullTime

		if err := rows.Scan(&areaID, &areaName, &sensorNo, &tempValue, &tempTS, &rhValue, &rhTS); err != nil {
			log.Printf("Error saat memindai baris data sensor gabungan: %v", err)
			continue
		}

		status := models.CombinedSensorStatus{
			AreaID:   areaID,
			AreaName: areaName,
			SensorNo: sensorNo,
		}

		tempKey := fmt.Sprintf("temp-%d-%d", areaID, sensorNo)
		tempOpStatus, tempFound := internalmodels.GetSensorOperationalStatus(tempKey)

		isOfflineByWorker := tempFound && tempOpStatus.IsOffline
		isOfflineByStaleData := tempTS.Valid && now.Sub(tempTS.Time) > offlineDuration

		if isOfflineByWorker || isOfflineByStaleData {
			status.Temp = &models.SensorValueStatus{
				Status: "offline",
			}
			if tempTS.Valid {
				status.Temp.Value = tempValue.Float64
				status.Temp.TS = tempTS.Time
			} else if tempFound {
				status.Temp.TS = tempOpStatus.LastSeen
			}
		} else if tempValue.Valid && tempTS.Valid {

			evalStatus := worker.EvaluateTemp(areaID, sensorNo, tempValue.Float64)
			status.Temp = &models.SensorValueStatus{
				Value:  tempValue.Float64,
				Status: mapSeverityToStatus(evalStatus.Severity),
				TS:     tempTS.Time,
			}
		}

		rhKey := fmt.Sprintf("rh-%d-%d", areaID, sensorNo)
		rhOpStatus, rhFound := internalmodels.GetSensorOperationalStatus(rhKey)

		isRhOfflineByWorker := rhFound && rhOpStatus.IsOffline
		isRhOfflineByStaleData := rhTS.Valid && now.Sub(rhTS.Time) > offlineDuration

		if isRhOfflineByWorker || isRhOfflineByStaleData {
			status.RH = &models.SensorValueStatus{
				Status: "offline",
			}
			if rhTS.Valid {
				status.RH.Value = rhValue.Float64
				status.RH.TS = rhTS.Time
			} else if rhFound {
				status.RH.TS = rhOpStatus.LastSeen
			}
		} else if rhValue.Valid && rhTS.Valid {
			evalStatus := worker.EvaluateRh(areaID, sensorNo, rhValue.Float64)
			status.RH = &models.SensorValueStatus{
				Value:  rhValue.Float64,
				Status: mapSeverityToStatus(evalStatus.Severity),
				TS:     rhTS.Time,
			}
		}
		statuses = append(statuses, status)
	}

	if statuses == nil {
		statuses = []models.CombinedSensorStatus{}
	}

	return c.Status(fiber.StatusOK).JSON(statuses)
}
