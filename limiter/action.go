package limiter

import "time"

type Action struct {
	Type      string
	Timestamp time.Time
}
