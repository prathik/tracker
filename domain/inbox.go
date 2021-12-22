package domain

import "time"

type Inbox struct {
	repo          InboxRepo
	CapturedItems []*InboxItem
}

func (i *Inbox) Load() {
	i.CapturedItems, _ = i.repo.GetAllInbox()
}

func NewInbox(repo InboxRepo) *Inbox {
	inbox := &Inbox{repo: repo}
	inbox.Load()
	return inbox
}

type InboxItem struct {
	CapturedTime time.Time
	Text string
	repo InboxRepo
}

func (i *InboxItem) Save() error {
	err := i.repo.Store(i)
	return err
}

func NewInboxItem(capturedTime time.Time, text string, repo InboxRepo) *InboxItem {
	return &InboxItem{CapturedTime: capturedTime, Text: text, repo: repo}
}