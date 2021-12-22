package domain

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestChallengeValues(t *testing.T) {
	type test struct {
		challenge string
		saveTime int
		wantErr  bool
	}

	tests := []test{
		{"PERFECT", 1,false},
		{"OVER", 1,false},
		{"UNDER", 1,false},
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