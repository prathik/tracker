package domain

type InboxRepo interface {
	Store(item *InboxItem) error
	GetAllInbox() ([]*InboxItem, error)
}