package models


type RawSensorData map[string]interface{}



type SensorPayload struct {
	AreaID *int     `json:"area"`
	DoorID *int     `json:"door"`
	No     *int     `json:"no"`
	TS     *string  `json:"ts"`
	Temp   *float64 `json:"temp"`
	RH     *float64 `json:"rh"`
	Prox   *int     `json:"prox"`
}


type SensorDataInput struct {
	CK     *int    `json:"ck"`
	AreaID *int    `json:"area"`
	DoorID *int    `json:"door"`
	TS     *string `json:"ts"`
}
