package cerrors_test

import (
	"fmt"
	"github.com/gocopper/copper/cerrors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithoutTags(t *testing.T) {
	err := fmt.Errorf("test-error-0; %w", cerrors.New(cerrors.New(nil, "test-error-2", map[string]interface{}{
		"tag": "val",
	}), "test-error-1", map[string]interface{}{
		"tag": "val",
	}))

	out := cerrors.WithoutTags(err)

	assert.NotContains(t, out.Error(), "tag")
}
