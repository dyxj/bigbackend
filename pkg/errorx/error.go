package errorx

import (
	"strings"
)

const validationErrorPrefix = "validation error"
const uniqueViolationErrorPrefix = "unique violation error"

type ValidationError struct {
	Properties map[string]string
}

func (e *ValidationError) Error() string {
	return writeErrorWithProperties(validationErrorPrefix, e.Properties)
}

type UniqueViolationError struct {
	Properties map[string]string
}

func (e *UniqueViolationError) Error() string {
	return writeErrorWithProperties(uniqueViolationErrorPrefix, e.Properties)
}

func writeErrorWithProperties(prefix string, properties map[string]string) string {
	if properties == nil || len(properties) == 0 {
		return prefix
	}

	sb := strings.Builder{}
	sb.WriteString(prefix)
	sb.WriteString(":")
	for k, v := range properties {
		sb.WriteString(" ")
		sb.WriteString(k)
		sb.WriteString(":")
		sb.WriteString(v)
	}
	return sb.String()
}
