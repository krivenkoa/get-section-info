package api

// Response is a generic response structure.
type Response struct {
	Status bool   `json:"status" example:"false"`
	Error  string `json:"error" example:"server doesn't respond'"`
}
