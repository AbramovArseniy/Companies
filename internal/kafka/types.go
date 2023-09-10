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

/*{: "some_alert_type", "uuid": "32dd6d6e-d7ae-4e4d-9b8e-8ba309fb4bbb", "ts": 1691897068000, "severity": "low", : "error|restore" }*/
