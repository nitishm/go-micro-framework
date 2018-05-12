package server

type Message struct {
	Service        Service `json:"service,omitempty"`
	ServiceMessage string  `json:"service_message,omitempty"`
}
