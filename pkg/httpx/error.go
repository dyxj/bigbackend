package httpx

const TimeoutResponseBody = `{error: "request timeout"}`

type ErrorResponse struct {
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}
