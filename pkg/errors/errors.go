package errors

import (
	pkgerrors "github.com/pkg/errors"
)

const (
	maskPatternCompileErrMsg = "compile mask pattern error"
	internalErrMsg           = "internal error"
)

var (
	MaskPatternCompileError = maskPatternError{}
	InternalError           = internalError{}
)

type maskPatternError struct{}

type internalError struct{}

func (maskPatternError) Error() string { return maskPatternCompileErrMsg }

func (internalError) Error() string { return internalErrMsg }

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

func ErrorIs(err error, target error) bool {
	return pkgerrors.Is(err, target)
}

func WithMessage(err error, message string) error {
	return pkgerrors.WithMessage(err, message)
}

func WithMessagef(err error, format string, args ...interface{}) error {
	return pkgerrors.WithMessagef(err, format, args...)
}

func ErrMaskPatternCompile(err error) error {
	return wrapErr(MaskPatternCompileError, err)
}

func ErrMaskPatternCompilef(err error, format string, args ...interface{}) error {
	return wrapErrf(MaskPatternCompileError, err, format, args...)
}

func ErrMaskPatternCompileRaw(message string) error {
	return wrapErr(MaskPatternCompileError, pkgerrors.New(message))
}

func ErrInternal(err error) error {
	return wrapErr(InternalError, err)
}

func ErrInternalf(err error, format string, args ...interface{}) error {
	return wrapErrf(InternalError, err, format, args...)
}

func ErrInternalRaw(message string) error {
	return wrapErr(InternalError, pkgerrors.New(message))
}
