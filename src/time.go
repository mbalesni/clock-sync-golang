package src

import (
	"fmt"
	"time"
)

type Time struct {
	Hours   int
	Minutes int
}

func (t *Time) Add(other *Time) *Time {

	minuteSum := t.Minutes + other.Minutes
	hourSum := t.Hours + other.Hours

	if minuteSum > 59 {

		t.Minutes = minuteSum - 60
		t.Hours = t.Hours + 1

	} else {

		t.Minutes = minuteSum

	}

	hourSum = t.Hours + other.Hours
	fmt.Println(hourSum, t, other)
	if hourSum > 24 {
		t.Hours = hourSum - 24
	} else {
		if t.Hours+other.Hours == 24 {

			t.Hours = 0

		} else {

			t.Hours = t.Hours + other.Hours

		}
	}

	return t

}

func (t *Time) Distance(other *Time) *Time {

	if t.Minutes > other.Minutes {

		return &Time{
			Minutes: MinutesDistance(t.Minutes, other.Minutes),
			Hours:   HoursDistance(t.Hours, other.Hours) - 1,
		}

	} else {
		return &Time{
			Minutes: MinutesDistance(t.Minutes, other.Minutes),
			Hours:   HoursDistance(t.Hours, other.Hours),
		}
	}

}

func CurrentTime() *Time {

	hour, min, _ := time.Now().Clock()

	return &Time{
		Minutes: min,
		Hours:   hour,
	}

}
