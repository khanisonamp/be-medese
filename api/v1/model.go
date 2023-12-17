package v1

type Result struct {
	Status    int         `json:"status,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	Message   interface{} `json:"message,omitempty"`
	MessageTh interface{} `json:"message_th,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}
