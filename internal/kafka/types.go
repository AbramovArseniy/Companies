package kafka

type TagChange struct {
	UUID      string  `json:"uuid"`
	TimeStamp float64 `json:"ts"`
	Value     float64 `json:"value"`
}
