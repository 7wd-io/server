package data

import "time"

func Now() time.Time {
	return time.Date(
		2024,
		time.January,
		1,
		10,
		0,
		0,
		0,
		time.UTC,
	)
}
