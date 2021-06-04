package limiter

import "time"

type Action struct {
	Type      string
	Timestamp time.Time
}

func NewAction(t string) *Action {
	a := new(Action)
	a.Timestamp = time.Now()
	a.Type = t
	return a
}
