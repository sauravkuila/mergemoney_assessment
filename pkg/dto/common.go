package dto

type CommonResponse struct {
	Status      bool     `json:"status"`
	Description string   `json:"description,omitempty"`
	Error       []string `json:"error,omitempty"`
}
