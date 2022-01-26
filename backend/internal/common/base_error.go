package common

type ErrorCode int

const (
	Unknown             ErrorCode = 0
	NotFound            ErrorCode = 1
	AlreadyExists       ErrorCode = 2
	FailToUnmarshalJson ErrorCode = 3
)

type WrappedError struct {
	Code  ErrorCode
	Error error
}
