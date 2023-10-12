package loggingdrain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMasking(t *testing.T) {
	t.Run("test number masking", func(t *testing.T) {
		beforeMaskStr := "D9 test 999, 888 1A ccc 3"
		expectMaskedStr := "D9 test <!NUM!>, <!NUM!> 1A ccc <!NUM!>"
		logMasker, _ := newLogMasker("<!", "!>")
		if err := logMasker.addInstruction("NUM", `\b\d+\b`); err != nil {
			t.Fatalf("add instruction error %s", err)
		}
		res := logMasker.mask(beforeMaskStr)
		assert.Equal(t, expectMaskedStr, res)
	})
	t.Run("test ip masking", func(t *testing.T) {
		beforeMaskStr := "the income ip is 127.0.0.1, abc, 10.3.24.13, 456"
		expectMaskedStr := "the income ip is <!IP!>, abc, <!IP!>, 456"
		logMasker, _ := newLogMasker("<!", "!>")
		if err := logMasker.addInstruction("IP", `\b(?:\d{1,3}\.){3}\d{1,3}\b`); err != nil {
			t.Fatalf("add instruction error %s", err)
		}
		res := logMasker.mask(beforeMaskStr)
		assert.Equal(t, expectMaskedStr, res)
	})
}
