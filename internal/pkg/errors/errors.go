package errors

type AppError struct {
	Code       string
	File       string
	Line       int
	Message    string
	StatusCode int
	Err        error
}

func (e *AppError) Error() string {
	return e.Message
}
