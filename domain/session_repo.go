package domain

import "time"

type SessionRepo interface {
	// Save a session
	Save(save *Session)

	// Query pulls the session items from daysBack number of days to current date inclusive.
	Query(duration time.Duration) ([]*Session, error)
	
	// Pop the last session added
	Pop()
}
