package domain

import (
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockSessionRepo(ctrl)
	sessionService := NewSessionService(repo)
	sessions := []*Session{{Time: time.Now(), Challenge: "flow"},
		{Time: time.Now().Add(-24 * time.Hour), Challenge: "anxiety"}}
	repo.EXPECT().Query(gomock.Any()).Return(sessions, nil)
	data, err := sessionService.ReportForPreviousDays(-1 * time.Hour)
	if err != nil {
		t.Fail()
	}

	if len(data) != 2 {
		t.Fail()
	}

	day1 := data[0]
	if day1.Sessions[0] != sessions[0] {
		t.Fail()
	}

	if day1.Count != 1 {
		t.Fail()
	}

	if day1.Day != sessions[0].Time.Format("2006-01-02") {
		t.Errorf("date mismatch")
	}

	day2 := data[1]
	if day2.Sessions[0] != sessions[1] {
		t.Fail()
	}

	if day2.Day != sessions[1].Time.Format("2006-01-02") {
		t.Errorf("date mismatch")
	}
}

func TestChallengeValues(t *testing.T) {
	type test struct {
		challenge string
		saveTime int
		wantErr  bool
	}

	tests := []test{
		{"flow", 1,false},
		{"anxiety", 1,false},
		{"boredom", 1,false},
		{"TEST", 0,true},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, tc := range tests {
		repo := NewMockSessionRepo(ctrl)
		sessionService := NewSessionService(repo)
		repo.EXPECT().Save(gomock.Any()).Times(tc.saveTime)
		err := sessionService.Save(&Session{Challenge: tc.challenge})
		if (err == nil) == tc.wantErr {
			t.Fail()
		}
	}
}
