package domain

import "time"

type SessionService struct {
	repo SessionRepo
}

func (s *SessionService) Save(item *Session) {
	s.repo.Save(item)
}

func (s *SessionService) QueryData(daysBack time.Duration) *Days {
	return s.repo.QueryData(daysBack)
}

func (s *SessionService) Pop() {
	s.repo.Pop()
}

func NewSessionService(repo SessionRepo) *SessionService {
	return &SessionService{repo: repo}
}