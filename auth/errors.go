package auth

type Error string

const (
	ErrInsufficientPermissions Error = "insufficent permissions" // TODO: 403
)

func (e Error) Error() string {
	return string(e)
}

func (e Error) Code() int {
	return 500
}
