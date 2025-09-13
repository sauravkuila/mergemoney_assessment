package dto

type AuthResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Status  bool        `json:"success"`
	Errors  []string    `json:"errors,omitempty"`
	Message string      `json:"description,omitempty"`
}
