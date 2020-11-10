package utils

type CancelledError struct{}

func (e *CancelledError) Error() string {
	return "cancelled by user"
}

func NewCancelledError() error {
	return &CancelledError{}
}
