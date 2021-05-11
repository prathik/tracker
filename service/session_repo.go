package service

type SessionRepo interface {
	Create(item *Item)
	GetWeekData() *WeekData
	// Pop the last item added
	Pop()
}
