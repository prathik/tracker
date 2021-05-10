package service

type SessionRepo interface {
	Create(item *Item)
	GetWeekData() *WeekData
}
