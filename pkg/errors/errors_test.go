package errors

import (
	stderrors "errors"
	"testing"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func createPatternCompileWrapErr() error {
	return ErrMaskPatternCompile(stderrors.New("mask pattern compile error"))
}

func TestErrors(t *testing.T) {
	t.Run("compile mask pattern error test", func(t *testing.T) {
		err1 := stderrors.New("compile mask pattern error")
		newErr := ErrMaskPatternCompile(err1)
		assert.Equal(t, "compile mask pattern error: compile mask pattern error", newErr.Error())
	})
	t.Run("mask pattern compile error raw error test", func(t *testing.T) {
		err1 := stderrors.New("compile mask pattern error")
		newErr := ErrMaskPatternCompilef(err1, "err1 error %s", "test")
		assert.Equal(t, "compile mask pattern error: err1 error test: compile mask pattern error", newErr.Error())
	})
	t.Run("nested cause test", func(t *testing.T) {
		err := createPatternCompileWrapErr()
		switch pkgerrors.Cause(err).(type) {
		case maskPatternError:
		default:
			t.Error("expect to equal InvalidError type")
		}
	})
	t.Run("nested compile mask pattern error test", func(t *testing.T) {
		err := createPatternCompileWrapErr()
		if !ErrorIs(err, MaskPatternCompileError) {
			t.Error("not equal")
		}
	})
}
