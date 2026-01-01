package idempotency

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithValue(t *testing.T) {
	ctx := context.Background()

	ctxMod := WithValue(ctx, "ikey")

	assert.Equal(t, "", FromContext(ctx))
	assert.Equal(t, "ikey", FromContext(ctxMod))
}
