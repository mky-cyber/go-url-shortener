package http

// URLResponse for /shortened endpoint response
type URLResponse struct {
	Result  string `json:"result,omitempty"`
	Message string `json:"message,omitempty"`
}
