package service

import "time"

// DayDataCollection represents the work done in the week.
type DayDataCollection struct {
	DayDataCollection []*DayData
}

// DayData is a collection of work done in a day.
type DayData struct {
	Time     time.Time
	Count    int
	WorkItem []*Item
}

// Item represents a single work item.
type Item struct {
	Time   time.Time // Time when the item was done
	Work   string    // Work that was done
	Joy    int       // Joy in the Work
	Impact int       // Business Impact of the Work
	Notes  string    // Any learnings from the Work is captured here
}
