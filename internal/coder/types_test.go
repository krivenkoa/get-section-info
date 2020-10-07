package coder

// TestRequest
type TestRequest struct {
	ID    string `json:"id" schema:"id"`
	Value string `json:"value" schema:"value"`
}

// TestResponse
type TestResponse struct {
	ID    string `json:"id" schema:"id"`
	Value string `json:"value" schema:"value"`
}
