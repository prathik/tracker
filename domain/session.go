package domain

import "time"

// Days represents the work done in multiple days.
type Days []*Day

// Day is a collection of sessions done in a day.
type Day struct {
	Day     string
	Count    int
	Sessions []*Session
}

// Session represents a single work session.
type Session struct {
	Time   time.Time // Time when the item was done
	Challenge string // Challenge is OVER, PERFECT or UNDER
}
