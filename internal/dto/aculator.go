package dto

import "time"

type AculatorMessage struct {
	ProcessId   string    `json:"process_id"`
	RoomId      string    `json:"room_id"`
	SensorId    string    `json:"sensor_id"`
	Service     string    `json:"service"`
	Value       float64   `json:"value"`
	PublishedOn time.Time `json:"published_on"`
}
