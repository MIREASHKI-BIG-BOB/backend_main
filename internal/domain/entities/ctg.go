package entities

type CTGData struct {
	SensorID     string  `json:"sensorID"`
	SecFromStart float32 `json:"secFromStart"`
	BPMChild     float32 `json:"BPMChild"`
	Uterus       float32 `json:"uterus"`
	Spasms       float32 `json:"spasms"`
}
