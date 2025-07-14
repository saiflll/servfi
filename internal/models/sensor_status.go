package models

import (
	"fmt"
	"log"

	//"strings"
	"sync"
	"time"

	"IoTT/internal/config"
	"IoTT/internal/database"
	"IoTT/internal/telegram"
)

type SensorOperationalStatus struct {
	SensorKey                   string
	SensorType                  string
	AreaID                      int
	SensorNo                    int
	DoorID                      int
	LastSeen                    time.Time
	IsOffline                   bool
	LastOfflineNotificationTime time.Time
}

var (
	sensorStatusRegistry      = make(map[string]*SensorOperationalStatus)
	sensorStatusRegistryMutex = &sync.Mutex{}
)

const OfflineReminderInterval = 1 * time.Hour

const (
	SensorTypeTemp = "temp"
	SensorTypeRH   = "rh"
	SensorTypeProx = "prox"
)

func RegisterOrUpdateSensorStatus(sensorKey, sensorType string, areaID, sensorNo, doorID int, timestamp time.Time) {
	sensorStatusRegistryMutex.Lock()
	defer sensorStatusRegistryMutex.Unlock()

	status, exists := sensorStatusRegistry[sensorKey]
	if !exists {
		status = &SensorOperationalStatus{
			SensorKey:  sensorKey,
			SensorType: sensorType,
			AreaID:     areaID,
			SensorNo:   sensorNo,
			DoorID:     doorID,
			LastSeen:   timestamp,
			IsOffline:  false,
		}
		sensorStatusRegistry[sensorKey] = status
		//log.Printf("Sensor %s terdaftar. Terakhir terlihat: %v", sensorKey, timestamp.Local())
	} else {
		status.LastSeen = timestamp // Selalu update waktu terakhir terlihat
		if status.IsOffline {
			//log.Printf("Sensor %s (%s) kembali online. Data baru pada: %v", getSensorFriendlyName(status), sensorKey, timestamp.Local().Format(time.RFC1123))
			onlineMsg := fmt.Sprintf("✅ **Sensor Online**\nSensor %s (%s) kembali mengirimkan data.\nData terakhir pada: %s", getSensorFriendlyName(status), sensorKey, timestamp.Local().Format(time.RFC1123))
			telegram.SendAlert(onlineMsg)
			status.IsOffline = false

		}
	}
}

func CheckAndNotifyOfflineSensors() {
	sensorStatusRegistryMutex.Lock()
	defer sensorStatusRegistryMutex.Unlock()
	now := time.Now().UTC()
	offlineThresholdDuration := time.Duration(config.PROX_THRESHOLDS.MaxIdleMinutes) * time.Minute

	for key, status := range sensorStatusRegistry {
		currentActualOfflineDuration := now.Sub(status.LastSeen)

		if currentActualOfflineDuration > offlineThresholdDuration { // Sensor telah melewati ambang batas offline
			if !status.IsOffline { // Ini adalah transisi pertama ke status offline (sejak terakhir online/terdaftar)
				actualOfflineMinutes := int(currentActualOfflineDuration.Minutes())
				//log.Printf("Sensor %s (%s) terdeteksi offline. Durasi: %d menit. Terakhir terlihat: %v", getSensorFriendlyName(status), key, actualOfflineMinutes, status.LastSeen.Local().Format(time.RFC1123))

				status.IsOffline = true
				status.LastOfflineNotificationTime = now // Catat waktu notifikasi pertama

				offlineMsg := fmt.Sprintf("⚠️ **Sensor Offline**\nSensor %s (%s) tidak mengirimkan data selama %d menit.\nTerakhir terlihat: %s\n🪛 Laporkan : [📞 Call . . . ](https://wa.me/+6282221294931)",
					getSensorFriendlyName(status), key, actualOfflineMinutes, status.LastSeen.Local().Format(time.RFC1123))

				telegram.SendAlert(offlineMsg)
			} else { // Sensor sudah dalam status offline, periksa untuk pengingat
				if now.Sub(status.LastOfflineNotificationTime) >= OfflineReminderInterval {
					totalOfflineMinutes := int(currentActualOfflineDuration.Minutes())
					//log.Printf("Mengirim pengingat offline untuk sensor %s (%s). Total durasi offline: %d menit. Notifikasi sebelumnya: %v", getSensorFriendlyName(status), key, totalOfflineMinutes, status.LastOfflineNotificationTime.Local().Format(time.RFC1123))

					reminderMsg := fmt.Sprintf("🕒 **Sensor Masih Offline (Pengingat)**\nSensor %s (%s) masih tidak mengirimkan data.\nTotal durasi offline: %d menit.\nNotifikasi terakhir dikirim %s.\nTerakhir terlihat: %s\n🪛 Laporkan : [📞 Call . . . ](https://facebook.com)",
						getSensorFriendlyName(status), key, totalOfflineMinutes, status.LastOfflineNotificationTime.Local().Format(time.RFC1123), status.LastSeen.Local().Format(time.RFC1123))
					telegram.SendAlert(reminderMsg)
					status.LastOfflineNotificationTime = now // Perbarui waktu untuk pengingat berikutnya
				}
			}
		}
	}
}

func getSensorFriendlyName(status *SensorOperationalStatus) string {

	switch status.SensorType {
	case "temp":
		areaName := database.GetAreaName(status.AreaID)
		return fmt.Sprintf("Suhu di Area %d Titik %d (%s ) ", status.AreaID, status.SensorNo, areaName)
	case "rh":
		areaName := database.GetAreaName(status.AreaID)
		return fmt.Sprintf("Kelembaban di Area %d Titik %d (%s)", status.AreaID, status.SensorNo, areaName)
	case "prox":
		doorName, _, associatedAreaName := database.GetDoorInfo(status.DoorID)
		return fmt.Sprintf("Proximity di Area %s Pintu %d (%s) ", associatedAreaName, status.DoorID, doorName)
	default:
		return status.SensorKey
	}
}

func StartOfflineDetectionWorker() {
	ticker := time.NewTicker(1 * time.Minute) // Periksa setiap 1 menit
	go func() {
		for {
			<-ticker.C // Tunggu tick berikutnya
			//log.Println("Inspection")
			CheckAndNotifyOfflineSensors()
		}
	}()
	log.Println("Offline Scan Begin . . . .")
}

// GetSensorOperationalStatus retrieves a copy of the operational status for a given sensor key.
// It returns the status and a boolean indicating if the sensor was found in the registry.
// It is safe for concurrent use.
func GetSensorOperationalStatus(sensorKey string) (SensorOperationalStatus, bool) {
	sensorStatusRegistryMutex.Lock()
	defer sensorStatusRegistryMutex.Unlock()

	status, exists := sensorStatusRegistry[sensorKey]
	if !exists {
		return SensorOperationalStatus{}, false
	}

	// Return a copy to prevent race conditions from the caller modifying the returned struct.
	return *status, true
}
