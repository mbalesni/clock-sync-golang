package src

import "time"

type Time struct {
	Hours   int
	Minutes int
}

func (t *Time) Add(other *Time) *Time {

	minuteSum := t.Minutes + other.Minutes
	hourSum := t.Hours + other.Hours

	if minuteSum > 59 {

		if minuteSum < 120 {

			t.Minutes = minuteSum - 60
			t.Hours = t.Hours + 1

		} else {

			t.Hours = t.Hours + 2

		}

	} else {

		t.Minutes = minuteSum

	}

	if hourSum > 24 {

		if hourSum < 48 {

			t.Hours = 48 - hourSum
			t.Hours = t.Hours + 1

		} else {

			t.Hours = t.Hours + 2

		}

	} else {
		t.Hours = t.Hours + other.Hours
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
