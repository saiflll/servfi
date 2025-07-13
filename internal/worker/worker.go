package worker

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	// "github.com/lib/pq" // HAPUS: Tidak digunakan lagi

	"IoTT/internal/config"
)

// InsertTemp diubah untuk menggunakan placeholder '?' yang lebih umum.
func InsertTemp(db *sql.DB, areaID int, no int, value float64, ts string) error {
	// GANTI: Placeholder $1 menjadi ?
	_, err := db.Exec("INSERT INTO temp (value, area_id, no, ts) VALUES (?, ?, ?, ?)",
		value, areaID, no, ts)
	if err != nil {
		log.Printf("Error inserting temp data: %v", err)
		return err
	}
	return nil
}

type TempBatchData struct {
	Value  float64
	AreaID int
	No     int
	TS     string
}

// BatchInsertTemp diubah total untuk menggunakan transaksi standar, bukan pq.CopyIn.
func BatchInsertTemp(tx *sql.Tx, data []TempBatchData) error {
	if len(data) == 0 {
		return nil
	}
	// GANTI: Menggunakan statement INSERT standar
	stmt, err := tx.Prepare("INSERT INTO temp (value, area_id, no, ts) VALUES (?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close() // Pastikan statement ditutup

	for _, d := range data {
		if _, err := stmt.Exec(d.Value, d.AreaID, d.No, d.TS); err != nil {
			// Tidak perlu close manual di sini karena ada defer
			return fmt.Errorf("error executing statement for temp data: %w", err)
		}
	}

	return nil // Commit akan dilakukan oleh pemanggil
}

// InsertRh diubah untuk menggunakan placeholder '?'
func InsertRh(db *sql.DB, areaID int, no int, value float64, ts string) error {
	_, err := db.Exec("INSERT INTO rh (value, area_id, no, ts) VALUES (?, ?, ?, ?)",
		value, areaID, no, ts)
	if err != nil {
		log.Printf("Error inserting rh data: %v", err)
		return err
	}
	return nil
}

type RhBatchData struct {
	Value  float64
	AreaID int
	No     int
	TS     string
}

// BatchInsertRh diubah total untuk menggunakan transaksi standar.
func BatchInsertRh(tx *sql.Tx, data []RhBatchData) error {
	if len(data) == 0 {
		return nil
	}
	stmt, err := tx.Prepare("INSERT INTO rh (value, area_id, no, ts) VALUES (?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare copy in for rh: %w", err)
	}
	defer stmt.Close()

	for _, d := range data {
		if _, err := stmt.Exec(d.Value, d.AreaID, d.No, d.TS); err != nil {
			return fmt.Errorf("failed to exec copy data for rh: %w", err)
		}
	}

	return nil
}

// InsertProx diubah untuk menggunakan placeholder '?'
func InsertProx(db *sql.DB, doorID int, value int, ts string) error {
	_, err := db.Exec("INSERT INTO prox (value, door_id, ts) VALUES (?, ?, ?)",
		value, doorID, ts)
	if err != nil {
		log.Printf("Error inserting prox data: %v", err)
		return err
	}
	return nil
}

type ProxBatchData struct {
	Value  int
	DoorID int
	TS     string
}

// BatchInsertProx diubah total untuk menggunakan transaksi standar.
func BatchInsertProx(tx *sql.Tx, data []ProxBatchData) error {
	if len(data) == 0 {
		return nil
	}
	stmt, err := tx.Prepare("INSERT INTO prox (value, door_id, ts) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare copy in for prox: %w", err)
	}
	defer stmt.Close()

	for _, d := range data {
		if _, err := stmt.Exec(d.Value, d.DoorID, d.TS); err != nil {
			return fmt.Errorf("failed to exec copy data for prox: %w", err)
		}
	}

	return nil
}

// =========================================================================
// TIDAK ADA PERUBAHAN PADA LOGIKA EVALUASI DAN UTILITAS DI BAWAH INI
// =========================================================================

type SafetyStatus struct {
	IsAlert   bool
	Message   string
	Severity  string
	Threshold float64
}

func getAreaName(areaID int) string {
	for _, area := range config.Areas {
		if area.ID == areaID {
			return area.Name
		}
	}
	return "Lokasi Tidak Dikenal"
}

func getThresholdConfig(thresholds []config.ThresholdConfig, areaID int, sensorNo int) (*config.Threshold, bool) {
	for _, t := range thresholds {
		if t.AreaID == areaID && t.SensorNo == sensorNo {
			return &t.Config, true
		}
	}
	return nil, false
}

func getMessageConfig(areaID int, sensorType string) (*config.MessageConfig, bool) {
	for _, m := range config.MessageConfigs {
		if m.AreaID == areaID && m.SensorType == sensorType {
			return &m, true
		}
	}
	return nil, false
}

func EvaluateTemp(areaID int, sensorNo int, currentValue float64) SafetyStatus {
	status := SafetyStatus{IsAlert: false, Severity: "NORMAL"}
	sensorCfg, ok := getThresholdConfig(config.TempThresholds, areaID, sensorNo)
	if !ok {
		return status
	}
	msgCfg, _ := getMessageConfig(areaID, "temp")
	location := getAreaName(areaID)
	if sensorCfg.Type == "ambient" {
		return status
	}
	if sensorCfg.UpperCritical != nil && currentValue > sensorCfg.UpperCritical.Limit {
		status.IsAlert = true
		status.Severity = "KRITIS_ATAS"
		status.Threshold = sensorCfg.UpperCritical.Limit
		if msgCfg != nil && msgCfg.UpperCriticalMsg != "" {
			status.Message = fmt.Sprintf(msgCfg.UpperCriticalMsg, config.PbP, location, currentValue, config.KtT, status.Threshold, config.IPbK, location)
		}
		return status
	}
	if sensorCfg.UpperWarning != nil && currentValue > sensorCfg.UpperWarning.Limit {
		status.IsAlert = true
		status.Severity = "WASPADA_ATAS"
		status.Threshold = sensorCfg.UpperWarning.Limit
		if msgCfg != nil && msgCfg.UpperWarningMsg != "" {
			status.Message = fmt.Sprintf(msgCfg.UpperWarningMsg, config.PwP, location, currentValue, config.KmT, status.Threshold, config.IPbW, location)
		}
		return status
	}
	if sensorCfg.LowerCritical != nil && currentValue < sensorCfg.LowerCritical.Limit {
		status.IsAlert = true
		status.Severity = "KRITIS_BAWAH"
		status.Threshold = sensorCfg.LowerCritical.Limit
		if msgCfg != nil && msgCfg.LowerCriticalMsg != "" {
			status.Message = fmt.Sprintf(msgCfg.LowerCriticalMsg, config.PbP, location, currentValue, config.KtR, status.Threshold, config.IPbK, location)
		}
		return status
	}
	if sensorCfg.LowerWarning != nil && currentValue < sensorCfg.LowerWarning.Limit {
		status.IsAlert = true
		status.Severity = "WASPADA_BAWAH"
		status.Threshold = sensorCfg.LowerWarning.Limit
		if msgCfg != nil && msgCfg.LowerWarningMsg != "" {
			status.Message = fmt.Sprintf(msgCfg.LowerWarningMsg, config.PwP, location, currentValue, config.KmR, status.Threshold, config.IPbW, location)
		}
		return status
	}
	return status
}

func EvaluateRh(areaID int, sensorNo int, currentValue float64) SafetyStatus {
	status := SafetyStatus{IsAlert: false, Severity: "NORMAL"}
	sensorCfg, ok := getThresholdConfig(config.RhThresholds, areaID, sensorNo)
	if !ok {
		return status
	}
	msgCfg, _ := getMessageConfig(areaID, "rh")
	location := getAreaName(areaID)
	if sensorCfg.UpperCritical != nil && currentValue > sensorCfg.UpperCritical.Limit {
		status.IsAlert = true
		status.Severity = "KRITIS_ATAS_RH"
		status.Threshold = sensorCfg.UpperCritical.Limit
		if msgCfg != nil && msgCfg.UpperCriticalMsg != "" {
			status.Message = fmt.Sprintf(msgCfg.UpperCriticalMsg, config.PbP, location, currentValue, config.KtT, status.Threshold, config.IPbK, location)
		}
		return status
	}
	if sensorCfg.UpperWarning != nil && currentValue > sensorCfg.UpperWarning.Limit {
		status.IsAlert = true
		status.Severity = "WASPADA_ATAS_RH"
		status.Threshold = sensorCfg.UpperWarning.Limit
		if msgCfg != nil && msgCfg.UpperWarningMsg != "" {
			status.Message = fmt.Sprintf(msgCfg.UpperWarningMsg, config.PwP, location, currentValue, config.KmT, status.Threshold, config.IPbW, location)
		}
		return status
	}
	if sensorCfg.LowerCritical != nil && currentValue < sensorCfg.LowerCritical.Limit {
		status.IsAlert = true
		status.Severity = "KRITIS_BAWAH_RH"
		status.Threshold = sensorCfg.LowerCritical.Limit
		if msgCfg != nil && msgCfg.LowerCriticalMsg != "" {
			status.Message = fmt.Sprintf(msgCfg.LowerCriticalMsg, config.PbP, location, currentValue, config.KtR, status.Threshold, config.IPbK, location)
		}
		return status
	}
	if sensorCfg.LowerWarning != nil && currentValue < sensorCfg.LowerWarning.Limit {
		status.IsAlert = true
		status.Severity = "WASPADA_BAWAH_RH"
		status.Threshold = sensorCfg.LowerWarning.Limit
		if msgCfg != nil && msgCfg.LowerWarningMsg != "" {
			status.Message = fmt.Sprintf(msgCfg.LowerWarningMsg, config.PwP, location, currentValue, config.KmR, status.Threshold, config.IPbW, location)
		}
		return status
	}
	return status
}

func GetFloat(v interface{}) (float64, error) {
	switch i := v.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case string:
		return strconv.ParseFloat(i, 64)
	default:
		return 0, fmt.Errorf("cannot convert type %T to float64", v)
	}
}

func GetInt(v interface{}) (int, error) {
	switch i := v.(type) {
	case float64:
		return int(i), nil
	case int:
		return i, nil
	case string:
		return strconv.Atoi(i)
	default:
		return 0, fmt.Errorf("cannot convert type %T to int", v)
	}
}

func ParseSensorNo(key, prefix string) (int, error) {
	if key == prefix {
		return 1, nil
	}
	if strings.HasPrefix(key, prefix) {
		noStr := strings.TrimPrefix(key, prefix)
		no, err := strconv.Atoi(noStr)
		if err != nil {
			return 0, fmt.Errorf("invalid sensor number in key '%s': %w", key, err)
		}
		return no, nil
	}
	return 0, fmt.Errorf("key '%s' does not have prefix '%s'", key, prefix)
}
