package domain

// Score gives a 10 point when in flow, gives a 3 point when challenge is over
// the current skill level and 0 when challenge is under skill level.
func Score(challenge string) int {
	switch challenge {
	case "PERFECT":
		return 10
	case "OVER":
		return 3
	default:
		return 0
	}
}
