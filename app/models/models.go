package models

import (
	"time"
)


type Area struct {
	ID   int    `json:"area_id"`
	Name string `json:"name"`
	CkID int    `json:"ck_id"`
}


type LoginRequest struct {
	
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}


type LatestSensorData struct {
	No          int       `json:"no"`
	Temperature *float64  `json:"temperature,omitempty"`
	RH          *float64  `json:"rh,omitempty"`
	TS          time.Time `json:"ts"`
}


type LatestSensorDataAllAreas struct {
	AreaID      int       `json:"area_id"`
	AreaName    string    `json:"area"`
	No          int       `json:"no"`
	Temperature *float64  `json:"temp,omitempty"`
	RH          *float64  `json:"rh,omitempty"`
	TS          time.Time `json:"ts"`
}


type SensorDataPoint struct {
	No    int       `json:"no"`
	Value float64   `json:"value"`
	TS    time.Time `json:"ts"`
}


type SensorDataRangeResponse struct {
	Temp []SensorDataPoint `json:"temp"`
	RH   []SensorDataPoint `json:"rh"`
}


type AreaSummary struct {
	AvgTemp    *float64 `json:"avg_temp,omitempty"`
	AvgRH      *float64 `json:"avg_rh,omitempty"`
	MaxTemp    *float64 `json:"max_temp,omitempty"`
	MinTemp    *float64 `json:"min_temp,omitempty"`
	LastTemp   *float64 `json:"last_temp,omitempty"`
	LastRH     *float64 `json:"last_rh,omitempty"`
	MinSetTemp *float64 `json:"min_set_temp,omitempty"`
	MaxSetTemp *float64 `json:"max_set_temp,omitempty"`
	MinSetRH   *float64 `json:"min_set_rh,omitempty"`
	MaxSetRH   *float64 `json:"max_set_rh,omitempty"`
}


type DetailedAlert struct {
	AreaName    string    `json:"nama_area"`
	SensorNo    int       `json:"no_sensor"`
	SensorType  string    `json:"tipe_sensor"`
	ValueBefore *float64  `json:"nilai_sebelum,omitempty"`
	ValueAfter  float64   `json:"nilai_setelah"`
	Description string    `json:"keterangan"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"waktu"`
}


type Alert struct {
	Area      string    `json:"area"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	Threshold float64   `json:"threshold"`
	TS        time.Time `json:"ts"`
}


type CK struct {
	ID   int    `json:"ck_id"`
	Name string `json:"name"`
}


type SensorValueStatus struct {
	Value  float64   `json:"value"`
	Status string    `json:"status"`
	TS     time.Time `json:"ts"`
}


type CombinedSensorStatus struct {
	AreaID   int                `json:"area_id"`
	AreaName string             `json:"area_name"`
	SensorNo int                `json:"sensor_no"`
	Temp     *SensorValueStatus `json:"temp,omitempty"`
	RH       *SensorValueStatus `json:"rh,omitempty"`
}
