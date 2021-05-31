package limiter

import (
	"errors"
	"time"
)

/* The Limiter struct saves all interactions in a map of lists indexed by user id going back the time limit.
When checking if a user is allowed to perform an action, it traverses the list for that userid.
If an item is older than the time limit, remove it and don't count.
If it is in the limit, then count. If the amount of interactions is higher than the limit, return an error */
type Limiter struct {
	TimeLimit time.Duration
	RateLimit int
	Logs      map[string][]*Action
}

func (l *Limiter) LogInteraction(userid string, action string) {
	ac := &Action{
		Timestamp: time.Now(),
		Type:      action,
	}
	l.Logs[userid] = append(l.Logs[userid], ac)
}

/* CheckAllowed counts the amount of log entries for a given userid, making sure to delete and not count the expired ones.
Returns an error if the amount of log entries exceeds the ratelimit */
func (l *Limiter) CheckAllowed(userid string) error {
	counter := 0
	expiredEntries := make([]*Action, 0)
	for i := 0; i < len(l.Logs[userid]); i++ {
		/* If the timestamp plus the timelimit is happened before "Now" */
		if l.Logs[userid][i].Timestamp.Add(l.TimeLimit).Before(time.Now()) {
			expiredEntries = append(expiredEntries, l.Logs[userid][i])
			continue
		} else {
			counter++
			continue
		}
	}
	/* remove entries */
	for i := 0; i < len(expiredEntries); i++ {
		l.removeAction(userid, expiredEntries[i])
	}

	if counter > l.RateLimit {
		return errors.New("rate limit exceeded")
	} else {
		return nil
	}
}

func (l *Limiter) removeAction(userid string, a *Action) {
	index := 0
	for i := 0; i < len(l.Logs[userid]); i++ {
		if l.Logs[userid][i] == a {
			index = i
			break
		}
	}
	l.Logs[userid][index] = l.Logs[userid][len(l.Logs[userid])-1]
	l.Logs[userid] = l.Logs[userid][:len(l.Logs[userid])-1]
}
