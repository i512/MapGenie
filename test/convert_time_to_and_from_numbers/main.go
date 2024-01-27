package convert_time_to_and_from_numbers

import "time"

type A struct {
	TimeToInt   time.Time
	TimeToInt64 time.Time
}

type B struct {
	TimeToInt   int
	TimeToInt64 int64
}
