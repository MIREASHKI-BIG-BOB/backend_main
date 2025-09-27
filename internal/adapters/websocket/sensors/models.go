package sensors

type MessageData struct {
	SensorID     string     `json:"sensorID"`
	SecFromStart float32    `json:"secFromStart"`
	Data         SensorData `json:"data"`
}

type SensorData struct {
	BPMChild float32 `json:"BPMChild"`
	Uterus   float32 `json:"uterus"`
	Spasms   float32 `json:"spasms"`
}
