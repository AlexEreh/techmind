package handlers

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
