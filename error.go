package main

type Error struct {
	status  int
	message string
}

func NewError(status int, message string) *Error {
	return &Error{
		status:  status,
		message: message,
	}
}

func (m *Error) Code() int {
	return m.status
}

func (m *Error) Error() string {
	return m.message
}
