package config

import "time"

type Area struct {
	ID   int
	Name string
}

type Threshold struct {
	UpperCritical *ThresholdLevel
	UpperWarning  *ThresholdLevel
	LowerWarning  *ThresholdLevel
	LowerCritical *ThresholdLevel
	Type          string
}

type ThresholdLevel struct {
	Limit float64
}

type MessageConfig struct {
	AreaID           int
	SensorType       string
	UpperCriticalMsg string
	UpperWarningMsg  string
	LowerWarningMsg  string
	LowerCriticalMsg string
}

type ThresholdConfig struct {
	AreaID   int
	SensorNo int
	Config   Threshold
}

type ProxThresholdsConfig struct {
	MaxIdleMinutes int
}

var (
	Areas           []Area
	MessageConfigs  []MessageConfig
	TempThresholds  []ThresholdConfig
	RhThresholds    []ThresholdConfig
	PROX_THRESHOLDS ProxThresholdsConfig
	WIBLocation     *time.Location
	PbP             = "PERINGATAN BAHAYA"
	PwP             = "PERINGATAN"
	KtT             = "DI ATAS AMBANG BATAS"
	KmT             = "MENDEKATI AMBANG BATAS ATAS"
	KmR             = "MENDEKATI AMBANG BATAS BAWAH"
	KtR             = "DI BAWAH AMBANG BATAS"
	IPbK            = "Segera lakukan tindakan korektif."
	IPbW            = "Segera lakukan pemantauan."
	L               = "Lokasi: {{.Location}}"
)

func init() {
	WIBLocation = time.FixedZone("WIB", 7*60*60)

	Areas = []Area{
		{ID: 1, Name: "Repacking Meat & Pawn"}, {ID: 2, Name: "Meat/Pawn Storage"}, {ID: 3, Name: "Chili/Mushroom Storage"},
		{ID: 4, Name: "Frozen Chili & Mushroom"}, {ID: 5, Name: "Ambient WH"}, {ID: 6, Name: "Packing Storage"},
		{ID: 7, Name: "Intermediate Room"}, {ID: 8, Name: "Corridor Room"}, {ID: 9, Name: "Meat Processing"},
		{ID: 10, Name: "Chilled Room"}, {ID: 11, Name: "Metos Room"}, {ID: 12, Name: "Wheat Flour Storage"},
		{ID: 13, Name: "Line Production"}, {ID: 14, Name: "IQF"}, {ID: 15, Name: "Secondary Packing"},
		{ID: 16, Name: "Intermediate Packing"}, {ID: 17, Name: "Frozen Storage FG"}, {ID: 18, Name: "Warehouse Dry FG"},
		{ID: 19, Name: "Loading Platform"},
	}

	MessageConfigs = []MessageConfig{

		{AreaID: 1, SensorType: "temp",
			UpperCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! (Batas Atas: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},

		{AreaID: 2, SensorType: "temp",
			UpperCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Atas: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			UpperWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Maks: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Min: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Bawah: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},
		{AreaID: 3, SensorType: "temp",
			UpperCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Atas: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			UpperWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Maks: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Min: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Bawah: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},
		// AreaID: 4 (Hanya Kritis Atas & Bawah)
		{AreaID: 4, SensorType: "temp",
			UpperCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Atas: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Bawah: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},

		{AreaID: 7, SensorType: "temp",
			UpperCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Atas: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			UpperWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Maks: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},
		{AreaID: 8, SensorType: "temp",
			UpperCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Atas: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			UpperWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Maks: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},
		{AreaID: 9, SensorType: "temp",
			UpperCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Atas: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			UpperWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Maks: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Min: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Bawah: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},
		{AreaID: 10, SensorType: "temp",
			UpperCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Atas: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			UpperWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Maks: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Min: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Bawah: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},
		{AreaID: 11, SensorType: "temp",
			UpperCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Atas: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			UpperWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Maks: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Min: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Bawah: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},
		{AreaID: 12, SensorType: "temp",
			UpperCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Atas: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			UpperWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Maks: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Min: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Bawah: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},
		{AreaID: 13, SensorType: "temp",
			UpperCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Atas: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			UpperWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Maks: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Min: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Bawah: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},

		{AreaID: 14, SensorType: "temp",
			UpperCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Atas: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			UpperWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Maks: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Bawah: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},
		{AreaID: 15, SensorType: "temp",
			UpperCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Atas: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			UpperWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Maks: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Min: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Bawah: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},
		{AreaID: 16, SensorType: "temp",
			UpperCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Atas: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			UpperWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Maks: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},
		{AreaID: 17, SensorType: "temp",
			UpperCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Atas: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			UpperWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Maks: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},
		{AreaID: 19, SensorType: "temp",
			UpperCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Atas: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			UpperWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Maks: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerWarningMsg:  "(%s)\n%s Suhu: %.1f°C %s (Normal Min: %.0f°C).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerCriticalMsg: "(%s)\n%s Suhu: %.1f°C %s! Batas Kritis Bawah: %.0f°C.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},

		// =========================================================================
		//                          KONFIGURASI KELEMBABAN (RH)
		// =========================================================================

		// AreaID: 11, 12 (Tanpa Kritis Atas & Peringatan Bawah)
		{AreaID: 11, SensorType: "rh",
			UpperWarningMsg:  "(%s)\n%s Kelembaban: %.1f%% %s (Normal Maks: %.0f%%).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerCriticalMsg: "(%s)\n%s Kelembaban: %.1f%% %s! Batas Kritis Bawah: %.0f%%.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},
		{AreaID: 12, SensorType: "rh",
			UpperWarningMsg:  "(%s)\n%s Kelembaban: %.1f%% %s (Normal Maks: %.0f%%).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerCriticalMsg: "(%s)\n%s Kelembaban: %.1f%% %s! Batas Kritis Bawah: %.0f%%.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},

		{AreaID: 13, SensorType: "rh",
			UpperCriticalMsg: "(%s)\n%s Kelembaban: %.1f%% %s! Batas Kritis Atas: %.0f%%.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			UpperWarningMsg:  "(%s)\n%s Kelembaban: %.1f%% %s (Normal Maks: %.0f%%).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerWarningMsg:  "(%s)\n%s Kelembaban: %.1f%% %s (Normal Min: %.0f%%).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerCriticalMsg: "(%s)\n%s Kelembaban: %.1f%% %s! Batas Kritis Bawah: %.0f%%.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},
		{AreaID: 15, SensorType: "rh",
			UpperCriticalMsg: "(%s)\n%s Kelembaban: %.1f%% %s! Batas Kritis Atas: %.0f%%.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			UpperWarningMsg:  "(%s)\n%s Kelembaban: %.1f%% %s (Normal Maks: %.0f%%).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerWarningMsg:  "(%s)\n%s Kelembaban: %.1f%% %s (Normal Min: %.0f%%).\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
			LowerCriticalMsg: "(%s)\n%s Kelembaban: %.1f%% %s! Batas Kritis Bawah: %.0f%%.\nHubungi: [call](https://wa.me/+6282221294931)\n%s\n%s",
		},
	}

	TempThresholds = []ThresholdConfig{
		{AreaID: 1, SensorNo: 1, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 10}}},
		{AreaID: 1, SensorNo: 2, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 10}}},
		{AreaID: 1, SensorNo: 3, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 10}}},
		{AreaID: 2, SensorNo: 1, Config: Threshold{LowerCritical: &ThresholdLevel{Limit: -22}, LowerWarning: &ThresholdLevel{Limit: -21}, UpperWarning: &ThresholdLevel{Limit: -19}, UpperCritical: &ThresholdLevel{Limit: -18}}},
		{AreaID: 2, SensorNo: 2, Config: Threshold{LowerCritical: &ThresholdLevel{Limit: -22}, LowerWarning: &ThresholdLevel{Limit: -21}, UpperWarning: &ThresholdLevel{Limit: -19}, UpperCritical: &ThresholdLevel{Limit: -18}}},
		{AreaID: 2, SensorNo: 3, Config: Threshold{LowerCritical: &ThresholdLevel{Limit: -22}, LowerWarning: &ThresholdLevel{Limit: -21}, UpperWarning: &ThresholdLevel{Limit: -19}, UpperCritical: &ThresholdLevel{Limit: -18}}},
		{AreaID: 3, SensorNo: 1, Config: Threshold{LowerCritical: &ThresholdLevel{Limit: 20}, LowerWarning: &ThresholdLevel{Limit: 21}, UpperWarning: &ThresholdLevel{Limit: 23}, UpperCritical: &ThresholdLevel{Limit: 24}}},
		{AreaID: 4, SensorNo: 1, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: -18}, LowerCritical: &ThresholdLevel{Limit: -22}}},
		{AreaID: 4, SensorNo: 2, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: -18}, LowerCritical: &ThresholdLevel{Limit: -22}}},
		{AreaID: 5, SensorNo: 1, Config: Threshold{Type: "ambient"}},
		{AreaID: 6, SensorNo: 1, Config: Threshold{Type: "ambient"}},
		{AreaID: 18, SensorNo: 1, Config: Threshold{Type: "ambient"}},
		{AreaID: 7, SensorNo: 1, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 10}, UpperWarning: &ThresholdLevel{Limit: 8}}},
		{AreaID: 8, SensorNo: 1, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 10}, UpperWarning: &ThresholdLevel{Limit: 8}}},
		{AreaID: 9, SensorNo: 1, Config: Threshold{UpperWarning: &ThresholdLevel{Limit: -18}, UpperCritical: &ThresholdLevel{Limit: -19}, LowerWarning: &ThresholdLevel{Limit: -21}, LowerCritical: &ThresholdLevel{Limit: -22}}},
		{AreaID: 10, SensorNo: 1, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 0}, UpperWarning: &ThresholdLevel{Limit: -1}, LowerWarning: &ThresholdLevel{Limit: -3}, LowerCritical: &ThresholdLevel{Limit: -4}}},
		{AreaID: 11, SensorNo: 1, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 24}, UpperWarning: &ThresholdLevel{Limit: 23}, LowerWarning: &ThresholdLevel{Limit: 21}, LowerCritical: &ThresholdLevel{Limit: 20}}},
		{AreaID: 12, SensorNo: 1, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 24}, UpperWarning: &ThresholdLevel{Limit: 23}, LowerWarning: &ThresholdLevel{Limit: 21}, LowerCritical: &ThresholdLevel{Limit: 20}}},
		{AreaID: 13, SensorNo: 1, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 24}, UpperWarning: &ThresholdLevel{Limit: 23}, LowerWarning: &ThresholdLevel{Limit: 21}, LowerCritical: &ThresholdLevel{Limit: 20}}},
		{AreaID: 13, SensorNo: 2, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 24}, UpperWarning: &ThresholdLevel{Limit: 23}, LowerWarning: &ThresholdLevel{Limit: 21}, LowerCritical: &ThresholdLevel{Limit: 20}}},
		{AreaID: 13, SensorNo: 3, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 24}, UpperWarning: &ThresholdLevel{Limit: 23}, LowerWarning: &ThresholdLevel{Limit: 21}, LowerCritical: &ThresholdLevel{Limit: 20}}},
		{AreaID: 13, SensorNo: 4, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 24}, UpperWarning: &ThresholdLevel{Limit: 23}, LowerWarning: &ThresholdLevel{Limit: 21}, LowerCritical: &ThresholdLevel{Limit: 20}}},
		{AreaID: 14, SensorNo: 1, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: -40}, UpperWarning: &ThresholdLevel{Limit: -41}, LowerCritical: &ThresholdLevel{Limit: -44}}},
		{AreaID: 14, SensorNo: 2, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: -40}, UpperWarning: &ThresholdLevel{Limit: -41}, LowerCritical: &ThresholdLevel{Limit: -44}}},
		{AreaID: 14, SensorNo: 3, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: -40}, UpperWarning: &ThresholdLevel{Limit: -41}, LowerCritical: &ThresholdLevel{Limit: -44}}},
		{AreaID: 14, SensorNo: 4, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: -40}, UpperWarning: &ThresholdLevel{Limit: -41}, LowerCritical: &ThresholdLevel{Limit: -44}}},
		{AreaID: 14, SensorNo: 5, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: -40}, UpperWarning: &ThresholdLevel{Limit: -41}, LowerCritical: &ThresholdLevel{Limit: -44}}},
		{AreaID: 14, SensorNo: 6, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: -40}, UpperWarning: &ThresholdLevel{Limit: -41}, LowerCritical: &ThresholdLevel{Limit: -44}}},
		{AreaID: 14, SensorNo: 7, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: -40}, UpperWarning: &ThresholdLevel{Limit: -41}, LowerCritical: &ThresholdLevel{Limit: -44}}},
		{AreaID: 14, SensorNo: 8, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: -40}, UpperWarning: &ThresholdLevel{Limit: -41}, LowerCritical: &ThresholdLevel{Limit: -44}}},
		{AreaID: 14, SensorNo: 9, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: -40}, UpperWarning: &ThresholdLevel{Limit: -41}, LowerCritical: &ThresholdLevel{Limit: -44}}},
		{AreaID: 15, SensorNo: 1, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 24}, UpperWarning: &ThresholdLevel{Limit: 23}, LowerWarning: &ThresholdLevel{Limit: 21}, LowerCritical: &ThresholdLevel{Limit: 20}}},
		{AreaID: 16, SensorNo: 1, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 10}, UpperWarning: &ThresholdLevel{Limit: 8}}},
		{AreaID: 16, SensorNo: 2, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 10}, UpperWarning: &ThresholdLevel{Limit: 8}}},
		{AreaID: 17, SensorNo: 1, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 20}, UpperWarning: &ThresholdLevel{Limit: 18}}},
		{AreaID: 17, SensorNo: 2, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 20}, UpperWarning: &ThresholdLevel{Limit: 18}}},
		{AreaID: 17, SensorNo: 3, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 20}, UpperWarning: &ThresholdLevel{Limit: 18}}},
		{AreaID: 17, SensorNo: 4, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 20}, UpperWarning: &ThresholdLevel{Limit: 18}}},
		{AreaID: 17, SensorNo: 5, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 20}, UpperWarning: &ThresholdLevel{Limit: 18}}},
		{AreaID: 17, SensorNo: 6, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 20}, UpperWarning: &ThresholdLevel{Limit: 18}}},
		{AreaID: 19, SensorNo: 1, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 12}, UpperWarning: &ThresholdLevel{Limit: 11}, LowerWarning: &ThresholdLevel{Limit: 9}, LowerCritical: &ThresholdLevel{Limit: 8}}},
		{AreaID: 19, SensorNo: 2, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 12}, UpperWarning: &ThresholdLevel{Limit: 11}, LowerWarning: &ThresholdLevel{Limit: 9}, LowerCritical: &ThresholdLevel{Limit: 8}}},
	}

	RhThresholds = []ThresholdConfig{
		{AreaID: 11, SensorNo: 1, Config: Threshold{UpperWarning: &ThresholdLevel{Limit: 80}, LowerCritical: &ThresholdLevel{Limit: 70}}},
		{AreaID: 12, SensorNo: 1, Config: Threshold{UpperWarning: &ThresholdLevel{Limit: 80}, LowerCritical: &ThresholdLevel{Limit: 70}}},
		{AreaID: 13, SensorNo: 1, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 60}, UpperWarning: &ThresholdLevel{Limit: 55}, LowerWarning: &ThresholdLevel{Limit: 45}, LowerCritical: &ThresholdLevel{Limit: 40}}},
		{AreaID: 13, SensorNo: 2, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 60}, UpperWarning: &ThresholdLevel{Limit: 55}, LowerWarning: &ThresholdLevel{Limit: 45}, LowerCritical: &ThresholdLevel{Limit: 40}}},
		{AreaID: 13, SensorNo: 3, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 60}, UpperWarning: &ThresholdLevel{Limit: 55}, LowerWarning: &ThresholdLevel{Limit: 45}, LowerCritical: &ThresholdLevel{Limit: 40}}},
		{AreaID: 13, SensorNo: 4, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 60}, UpperWarning: &ThresholdLevel{Limit: 55}, LowerWarning: &ThresholdLevel{Limit: 45}, LowerCritical: &ThresholdLevel{Limit: 40}}},
		{AreaID: 15, SensorNo: 1, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 60}, UpperWarning: &ThresholdLevel{Limit: 55}, LowerWarning: &ThresholdLevel{Limit: 45}, LowerCritical: &ThresholdLevel{Limit: 40}}},
		{AreaID: 15, SensorNo: 2, Config: Threshold{UpperCritical: &ThresholdLevel{Limit: 60}, UpperWarning: &ThresholdLevel{Limit: 55}, LowerWarning: &ThresholdLevel{Limit: 45}, LowerCritical: &ThresholdLevel{Limit: 40}}},
	}

	PROX_THRESHOLDS = ProxThresholdsConfig{
		MaxIdleMinutes: 10,
	}
}
