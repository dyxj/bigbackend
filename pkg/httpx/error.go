package httpx

type ErrorResponse struct {
	Code    errorCode         `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

type errorCode string

const (
	CodeBadRequest               errorCode = "bad_request"
	CodeServerError              errorCode = "server_error"
	CodeEntityNotFound           errorCode = "entity_not_found"
	CodeDuplicateEntity          errorCode = "duplicate_entity"
	CodeIdempotencyKeyInProgress errorCode = "idempotency_key_in_progress"
)

func (e errorCode) String() string {
	return string(e)
}
