package domain

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func Test_formatDelta(t *testing.T) {
	type args struct {
		now    time.Time
		doneOn time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1m delta", args{now: time.Now(), doneOn: time.Now().Add(-1*time.Minute)}, "-1m"},
		{"10m delta", args{now: time.Now(), doneOn: time.Now().Add(-10*time.Minute)}, "-10m"},
		{"59m delta", args{now: time.Now(), doneOn: time.Now().Add(-59*time.Minute)}, "-59m"},
		{"hour delta", args{now: time.Now(), doneOn: time.Now().Add(-60*time.Minute)}, "-1h0m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatDelta(tt.args.now, tt.args.doneOn); got != tt.want {
				t.Errorf("formatDelta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHourAndMinuteFromString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name       string
		args       args
		wantHour   int
		wantMinute int
		wantErr    bool
	}{
		{"returns error when empty time", args{str: ""}, 0, 0, true},
		{"returns error when incorrect format", args{str: "sfsdfsd"}, 0, 0, true},
		{"returns error when two colons present", args{str: "10:10:10"}, 0, 0, true},
		{"returns error when no minute", args{str: "10:"}, 0, 0, true},
		{"returns error when no hour", args{str: ":10"}, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHour, gotMinute, err := HourAndMinuteFromString(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("HourAndMinuteFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHour != tt.wantHour {
				t.Errorf("HourAndMinuteFromString() gotHour = %v, want %v", gotHour, tt.wantHour)
			}
			if gotMinute != tt.wantMinute {
				t.Errorf("HourAndMinuteFromString() gotMinute = %v, want %v", gotMinute, tt.wantMinute)
			}
		})
	}
}

func TestPrintWeekData(t *testing.T) {
	t.Run("It should only give today's data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := NewMockSessionRepo(ctrl)
		sessionService := NewSessionService(repo)
		prev, _ := time.Parse("02-01-2006", "01-01-2022")
		cur, _ := time.Parse("02-01-2006", "02-01-2022")
		sessions := []*Session{{Time: cur, Challenge: "flow"},
			{Time: prev, Challenge: "anxiety"}}
		repo.EXPECT().Query(gomock.Any()).Return(sessions, nil)
		data, _ := PrintWeekData(sessionService)

		if len(data) != 2 {
			t.Fail()
		}
		fmt.Println(data)

		first := data[0]

		if first[0] != "2022-01-02" {
			t.Fail()
		}

		if first[1] != "1" {
			t.Fail()
		}

		second := data[1]

		if second[0] != "2022-01-01" {
			t.Fail()
		}

		if second[1] != "1" {
			t.Fail()
		}
	})

	t.Run("It should only give today's data in sorted order", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := NewMockSessionRepo(ctrl)
		sessionService := NewSessionService(repo)
		prev, _ := time.Parse("02-01-2006", "02-01-2022")
		cur, _ := time.Parse("02-01-2006", "01-01-2022")
		sessions := []*Session{{Time: cur, Challenge: "flow"},
			{Time: prev, Challenge: "anxiety"}}
		repo.EXPECT().Query(gomock.Any()).Return(sessions, nil)
		data, _ := PrintWeekData(sessionService)

		if len(data) != 2 {
			t.Fail()
		}
		fmt.Println(data)

		first := data[0]

		if first[0] != "2022-01-02" {
			t.Fail()
		}

		if first[1] != "1" {
			t.Fail()
		}

		second := data[1]

		if second[0] != "2022-01-01" {
			t.Fail()
		}

		if second[1] != "1" {
			t.Fail()
		}
	})
}