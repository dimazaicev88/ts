package dto

type ServerMetadata struct {
	NatsURL       string `json:"natsUrl"`
	StreamName    string `json:"streamName"`
	ConsumerName  string `json:"consumerName"`
	CentrifugeURL string `json:"centrifugeURL"`
}
