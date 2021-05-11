package service

type SessionService struct {
	repo SessionRepo
}

func (s *SessionService) Create(item *Item) {
	s.repo.Create(item)
}

func (s *SessionService) QueryData(daysBack int) *DayDataCollection {
	return s.repo.QueryData(daysBack)
}

func (s *SessionService) Pop() {
	s.repo.Pop()
}

func NewSessionService(repo SessionRepo) *SessionService {
	return &SessionService{repo: repo}
}