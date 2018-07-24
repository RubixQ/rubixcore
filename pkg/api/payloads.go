package api

// CounterRegPayload is sent by counters to register on server
type CounterRegPayload struct {
	CounterID string `json:"counterId"`
	QueueID   string `json:"queueId"`
}

// WSPayload wraps around ws responses
type WSPayload struct {
	PayloadType string      `json:"type"`
	Data        interface{} `json:"data"`
}
