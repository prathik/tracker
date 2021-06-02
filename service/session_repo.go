package service

import "time"

type SessionRepo interface {
	Create(item *Item)

	// QueryData pulls the data items from daysBack number of days to current date inclusive.
	QueryData(duration time.Duration) *DayDataCollection
	// Pop the last item added
	Pop()
}
