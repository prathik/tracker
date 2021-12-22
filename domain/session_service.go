package domain

import (
	"errors"
	"time"
)

type SessionService struct {
	repo SessionRepo
}

func validChallenge(challenge string) bool {
	switch challenge {
	case "PERFECT":
		return true
	case "OVER":
		return true
	case "UNDER":
		return true
	default:
		return false
	}
}

func (s *SessionService) Save(item *Session) error {
	if !validChallenge(item.Challenge) {
		return errors.New("invalid challenge level")
	}
	s.repo.Save(item)
	return nil
}

func (s *SessionService) QueryData(daysBack time.Duration) (*Days, error) {
	return s.repo.QueryData(daysBack)
}

func (s *SessionService) Pop() {
	s.repo.Pop()
}

func NewSessionService(repo SessionRepo) *SessionService {
	return &SessionService{repo: repo}
}