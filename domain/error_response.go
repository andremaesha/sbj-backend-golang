package domain

type ErrorResponse struct {
	RequestId string            `json:"request_id,omitempty"`
	Message   string            `json:"message"`
	Errors    map[string]string `json:"errors,omitempty"`
}
