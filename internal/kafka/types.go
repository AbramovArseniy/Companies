package kafka

type TagChange struct {
	UUID      string  `json:"uuid"`
	TimeStamp float64 `json:"ts"`
	Value     float64 `json:"value"`
}

type Alert struct {
	MsgKind   string  `json:"msg_kind"`
	Type      string  `json:"type"`
	TagID     string  `json:"uuid"`
	TimeStamp float64 `json:"ts"`
	Severity  string  `json:"severity"`
	State     string  `json:"alert_state"`
}
