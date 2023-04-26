package gpt

type ChatMessageTypeEnum int

const (
	ZERO ChatMessageTypeEnum = iota
	QUESTION
	ANSWER
)

func (c ChatMessageTypeEnum) String() string {
	switch c {
	case ZERO:
		return "ZERO"
	case QUESTION:
		return "QUESTION"
	case ANSWER:
		return "ANSWER"
	default:
		return "UNKNOWN"
	}
}
