package domain

import "time"

type SessionRepo interface {
	// Save a session
	Save(save *Session)

	// QueryData pulls the data items from daysBack number of days to current date inclusive.
	QueryData(duration time.Duration) (*Days, error)
	
	// Pop the last session added
	Pop()
}
