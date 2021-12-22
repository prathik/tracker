package domain

import "time"

// Days represents the work done in multiple days.
type Days struct {
	Days []*Day
}

// Day is a collection of sessions done in a day.
type Day struct {
	Time     time.Time
	Count    int
	Sessions []*Session
}

// Session represents a single work session.
type Session struct {
	Time   time.Time // Time when the item was done
	Joy    int       // Joy in the Work
	Impact int       // Business Impact of the Work
	Notes  string    // Any learnings from the Work is captured here
}
