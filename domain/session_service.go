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
	case "flow":
		return true
	case "anxiety":
		return true
	case "boredom":
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

func (s *SessionService) ReportForPreviousDays(daysBack time.Duration) (Days, error) {
	sessions, err := s.repo.Query(daysBack)

	if err != nil {
		return nil, err
	}

	// Group by date
	var daySessionsMap = make(map[string][]*Session)
	for _, session := range sessions {
		daySessionsMap[session.Time.Format("2006-01-02")]=
			append(daySessionsMap[session.Time.Format("2006-01-02")], session)
	}

	var days Days
	for day, daySessions := range daySessionsMap {
		days = append(days, &Day{Sessions: daySessions, Count: len(daySessions), Day: day})
	}
	return days, err
}

func (s *SessionService) Pop() {
	s.repo.Pop()
}

func NewSessionService(repo SessionRepo) *SessionService {
	return &SessionService{repo: repo}
}
