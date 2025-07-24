package cerr

type CustomError struct {
	Code    int
	Message string
	Err     error
}

func New(code int, message string, err error) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *CustomError) Error() string {
	return e.Message
}

func (e *CustomError) WithErr(err error) *CustomError {
	e.Err = err
	return e
}
