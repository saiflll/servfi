package models

import (
	"fmt"
	"log"
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

	} else {
		status.LastSeen = timestamp
		if status.IsOffline {

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

		if currentActualOfflineDuration > offlineThresholdDuration {
			if !status.IsOffline {
				actualOfflineMinutes := int(currentActualOfflineDuration.Minutes())

				status.IsOffline = true
				status.LastOfflineNotificationTime = now

				offlineMsg := fmt.Sprintf("⚠️ **Sensor Offline**\nSensor %s (%s) tidak mengirimkan data selama %d menit.\nTerakhir terlihat: %s\n🪛 Laporkan : [📞 Call . . . ](https://wa.me/+6282221294931)",
					getSensorFriendlyName(status), key, actualOfflineMinutes, status.LastSeen.Local().Format(time.RFC1123))

				telegram.SendAlert(offlineMsg)
			} else {
				if now.Sub(status.LastOfflineNotificationTime) >= OfflineReminderInterval {
					totalOfflineMinutes := int(currentActualOfflineDuration.Minutes())

					reminderMsg := fmt.Sprintf("🕒 **Sensor Masih Offline (Pengingat)**\nSensor %s (%s) masih tidak mengirimkan data.\nTotal durasi offline: %d menit.\nNotifikasi terakhir dikirim %s.\nTerakhir terlihat: %s\n🪛 Laporkan : [📞 Call . . . ](https://facebook.com)",
						getSensorFriendlyName(status), key, totalOfflineMinutes, status.LastOfflineNotificationTime.Local().Format(time.RFC1123), status.LastSeen.Local().Format(time.RFC1123))
					telegram.SendAlert(reminderMsg)
					status.LastOfflineNotificationTime = now
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
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			<-ticker.C

			CheckAndNotifyOfflineSensors()
		}
	}()
	log.Println("Offline Scan Begin . . . .")
}

func GetSensorOperationalStatus(sensorKey string) (SensorOperationalStatus, bool) {
	sensorStatusRegistryMutex.Lock()
	defer sensorStatusRegistryMutex.Unlock()

	status, exists := sensorStatusRegistry[sensorKey]
	if !exists {
		return SensorOperationalStatus{}, false
	}

	return *status, true
}
