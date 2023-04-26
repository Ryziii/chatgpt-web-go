package gpt

type ChatMessageStatusEnum int

const (
	INIT ChatMessageStatusEnum = iota
	PART_SUCCESS

	COMPLETE_SUCCESS

	ERROR
)

func (c ChatMessageStatusEnum) String() string {
	switch c {
	case INIT:
		return "INIT"
	case PART_SUCCESS:
		return "PART_SUCCESS"
	case COMPLETE_SUCCESS:
		return "COMPLETE_SUCCESS"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
