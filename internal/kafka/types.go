package kafka

import "time"

type TagChange struct {
	UUID      string    `json:"uuid"`
	TimeStamp time.Time `json:"ts"`
	Value     float64   `json:"value"`
}
