package models

// RawSensorData adalah tipe lama untuk JSON dinamis.
type RawSensorData map[string]interface{}

// SensorPayload adalah struct baru yang lebih terstruktur untuk menerima data sensor.
// Menggunakan pointer memungkinkan field menjadi opsional.
type SensorPayload struct {
	AreaID *int     `json:"area"`
	DoorID *int     `json:"door"`
	No     *int     `json:"no"`
	TS     *string  `json:"ts"`
	Temp   *float64 `json:"temp"`
	RH     *float64 `json:"rh"`
	Prox   *int     `json:"prox"`
}

// SensorDataInput tidak lagi digunakan secara langsung oleh handler utama.
type SensorDataInput struct {
	CK     *int    `json:"ck"`
	AreaID *int    `json:"area"`
	DoorID *int    `json:"door"`
	TS     *string `json:"ts"`
}
