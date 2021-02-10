package cerrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/v2/cerrors"
)

func TestNew(t *testing.T) {
	t.Parallel()

	err := cerrors.New(nil, "test-err", nil)

	cErr, ok := err.(cerrors.Error) //nolint:errorlint

	assert.True(t, ok)
	assert.Equal(t, "test-err", cErr.Message)
	assert.Nil(t, cErr.Cause)
	assert.Nil(t, cErr.Tags)
}

func TestWithTags_StdErr(t *testing.T) {
	t.Parallel()

	err := cerrors.WithTags(errors.New("test-err"), map[string]interface{}{ //nolint:goerr113
		"key": "val",
	})

	cErr, ok := err.(cerrors.Error) //nolint:errorlint

	assert.True(t, ok)
	assert.Equal(t, "test-err", cErr.Message)
	assert.Contains(t, cErr.Tags, "key")
	assert.Equal(t, cErr.Tags["key"], "val")
	assert.Nil(t, cErr.Cause)
}

func TestWithTags_CErr(t *testing.T) {
	t.Parallel()

	err := cerrors.WithTags(cerrors.New(nil, "test-cerr", nil), map[string]interface{}{
		"key": "val",
	})

	cErr, ok := err.(cerrors.Error) //nolint:errorlint

	assert.True(t, ok)
	assert.Equal(t, "test-cerr", cErr.Message)
	assert.Contains(t, cErr.Tags, "key")
	assert.Equal(t, cErr.Tags["key"], "val")
	assert.Nil(t, cErr.Cause)
}

func TestError_Unwrap(t *testing.T) {
	t.Parallel()

	err := cerrors.New(errors.New("cause-err"), "test-err", nil) //nolint:goerr113
	cause := errors.Unwrap(err)

	assert.NotNil(t, cause)
	assert.EqualError(t, cause, "cause-err")
}

func TestError_Unwrap_NoCause(t *testing.T) {
	t.Parallel()

	err := cerrors.New(nil, "test-err", nil)

	assert.Nil(t, errors.Unwrap(err))
}

func TestError_Error(t *testing.T) {
	t.Parallel()

	err := cerrors.New(nil, "test-err", nil)

	assert.NotNil(t, err)
	assert.Equal(t, "test-err", err.Error())
}

func TestError_Error_Cause(t *testing.T) {
	t.Parallel()

	err := cerrors.New(errors.New("cause-err"), "test-err", nil) //nolint:goerr113

	expectedErr := `test-err because
> cause-err`

	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err.Error())
}

func TestError_Error_Tags(t *testing.T) {
	t.Parallel()

	err := cerrors.New(nil, "test-err", map[string]interface{}{
		"key": "val",
	})

	assert.NotNil(t, err)
	assert.Equal(t, "test-err where key=val", err.Error())
}

func TestError_Error_PtrTags(t *testing.T) {
	t.Parallel()

	val := "val"
	err := cerrors.New(nil, "test-err", map[string]interface{}{
		"key": &val,
	})

	assert.NotNil(t, err)
	assert.Equal(t, "test-err where key=val", err.Error())
}
