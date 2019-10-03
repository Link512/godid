package godid

import "fmt"

//DidError is an error that occured due to invalid user data
type DidError struct {
	message string
}

func (e DidError) Error() string {
	return e.message
}

func didErrorf(format string, args ...interface{}) DidError {
	return DidError{
		message: fmt.Sprintf(format, args...),
	}
}
