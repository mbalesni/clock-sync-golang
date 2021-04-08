package src

import (
	"fmt"
	"testing"
)

func TestTime(t *testing.T) {

	twoAm := &Time{Hours: 2, Minutes: 59}

	satanTime := &Time{Hours: 23, Minutes: 0}

	twoAmSatanTimeDiff := twoAm.Distance(satanTime)

	if twoAmSatanTimeDiff.Hours != 20 && twoAmSatanTimeDiff.Minutes != 59 {

		t.Fail()

	}

	satanTimeTwoAmDiff := satanTime.Distance(twoAm)

	if satanTimeTwoAmDiff.Hours != 3 && satanTimeTwoAmDiff.Minutes != 1 {

		t.Fail()

	}

	fifteenMinutes := &Time{Hours: 0, Minutes: 15}
	twoFifteen := twoAm.Add(fifteenMinutes)

	if twoFifteen.Hours != 3 || twoFifteen.Minutes != 14 {
		fmt.Println(twoFifteen)
		t.Fail()
	}

}
