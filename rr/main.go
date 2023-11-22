package rr

func New(message string) AppError {
	return AppError{
		Message: message,
	}
}

type AppError struct {
	Message string `json:"err"`
}

func (dst AppError) Error() string {
	return dst.Message
}
