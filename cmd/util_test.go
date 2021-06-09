package cmd

import (
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
