package enum

type Status string

const (
	StatusInProgress Status = "in_progress"
	StatusClosed     Status = "closed"
)

func (s Status) String() string {
	return string(s)
}
