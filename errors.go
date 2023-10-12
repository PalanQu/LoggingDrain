package loggingdrain

import (
	pkgerrors "github.com/pkg/errors"
)

const (
	maskPatternCompileErrMsg = "compile mask pattern error"
	internalErrMsg           = "internal error"
)

var (
	maskPatternCompileError = MaskPatternError{}
	internalError           = InternalError{}
)

type MaskPatternError struct{}

type InternalError struct{}

func (MaskPatternError) Error() string { return maskPatternCompileErrMsg }

func (InternalError) Error() string { return internalErrMsg }

func wrapErr(wrapedErr error, err error) error {
	if err == nil {
		return pkgerrors.Wrap(wrapedErr, "")
	}
	return pkgerrors.Wrap(wrapedErr, err.Error())
}

func wrapErrf(
	wrapedErr error, err error, format string, args ...interface{}) error {
	e := pkgerrors.WithMessagef(wrapedErr, format, args...)
	if err == nil {
		return pkgerrors.Wrap(wrapedErr, e.Error())
	}
	return pkgerrors.Wrap(e, err.Error())
}

func errorIs(err error, target error) bool {
	return pkgerrors.Is(err, target)
}

func withMessage(err error, message string) error {
	return pkgerrors.WithMessage(err, message)
}

func withMessagef(err error, format string, args ...interface{}) error {
	return pkgerrors.WithMessagef(err, format, args...)
}

func errMaskPatternCompile(err error) error {
	return wrapErr(maskPatternCompileError, err)
}

func errMaskPatternCompilef(err error, format string, args ...interface{}) error {
	return wrapErrf(maskPatternCompileError, err, format, args...)
}

func errMaskPatternCompileRaw(message string) error {
	return wrapErr(maskPatternCompileError, pkgerrors.New(message))
}

func errInternal(err error) error {
	return wrapErr(internalError, err)
}

func errInternalf(err error, format string, args ...interface{}) error {
	return wrapErrf(internalError, err, format, args...)
}

func errInternalRaw(message string) error {
	return wrapErr(internalError, pkgerrors.New(message))
}
