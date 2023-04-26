package global

type SystemError struct {
	Msg string
}

func (e *SystemError) Error() string {
	return e.Msg
}
