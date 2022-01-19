package domain

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestPrintWeekData(t *testing.T) {
	t.Run("It should only give current week's data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := NewMockSessionRepo(ctrl)
		sessionService := NewSessionService(repo)
		prev, _ := time.Parse("02-01-2006", "03-01-2022")
		cur, _ := time.Parse("02-01-2006 15:04", "04-01-2022 10:00")
		sessions := []*Session{{Time: cur, Challenge: "flow"},
			{Time: cur.Add(1 * time.Hour), Challenge: "flow"},
			{Time: prev, Challenge: "anxiety"},
			{Time: prev.Add(1 * time.Hour), Challenge: "anxiety"}}

		duration, _ := time.ParseDuration("34h") // 7 days

		repo.EXPECT().Query(duration).Return(sessions, nil)
		data, _ := GenerateWeekReport(cur, sessionService)

		if data.FlowRatio != 0.5 {
			t.Errorf("invaild flow %f", data.FlowRatio)
		}

		if data.AverageCount != 2 {
			t.Errorf("invalid average")
		}

		if len(data.Raw) != 2 {
			t.Fail()
		}

		first := data.Raw[0]

		if first[0] != "2022-01-04" {
			t.Fail()
		}

		if first[1] != "2" {
			t.Fail()
		}

		second := data.Raw[1]

		if second[0] != "2022-01-03" {
			t.Fail()
		}

		if second[1] != "2" {
			t.Errorf("incorrect count for previous day")
		}
	})

	t.Run("It should only give current week's data in sorted order", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := NewMockSessionRepo(ctrl)
		sessionService := NewSessionService(repo)
		prev, _ := time.Parse("02-01-2006", "02-01-2022")
		cur, _ := time.Parse("02-01-2006", "01-01-2022")
		sessions := []*Session{{Time: cur, Challenge: "flow"},
			{Time: prev, Challenge: "anxiety"}}
		repo.EXPECT().Query(gomock.Any()).Return(sessions, nil)
		data, _ := GenerateWeekReport(time.Now(), sessionService)

		if len(data.Raw) != 2 {
			t.Fail()
		}
		fmt.Println(data)

		first := data.Raw[0]

		if first[0] != "2022-01-02" {
			t.Fail()
		}

		if first[1] != "1" {
			t.Fail()
		}

		second := data.Raw[1]

		if second[0] != "2022-01-01" {
			t.Fail()
		}

		if second[1] != "1" {
			t.Fail()
		}
	})
}

func Test_timeToStartOfTheWeek(t *testing.T) {
	type args struct {
		current time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{"monday 10th hour",
			args{current: time.Date(2022, 01, 17, 10, 0, 0, 0, time.UTC)},
			10 * time.Hour},
		{"sunday 10th hour",
			args{current: time.Date(2022, 01, 16, 10, 0, 0, 0, time.UTC)},
			154 * time.Hour},
		{"saturday 10th hour",
			args{current: time.Date(2022, 01, 15, 10, 0, 0, 0, time.UTC)},
			130 * time.Hour},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timeToStartOfTheWeek(tt.args.current); got != tt.want {
				t.Errorf("timeToStartOfTheWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}