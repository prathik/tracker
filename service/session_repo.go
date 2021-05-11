package service

type SessionRepo interface {
	Create(item *Item)

	// QueryData pulls the data items from daysBack number of days to current date inclusive.
	QueryData(daysBack int) *DayDataCollection
	// Pop the last item added
	Pop()
}
