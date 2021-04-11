package src

import "fmt"

func absInt(x int) int {

	if x < 0 {

		return -x

	}
	return x

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func mod(a, b int) int {
	m := a % b
	if a < 0 && b < 0 {
		m -= b
	}
	if a < 0 && b > 0 {
		m += b
	}
	return m
}

func unidirectionalRingDistance(from, to, maxSize int) int {

	toFrom := to - from

	if toFrom == 0 {
		return 0
	} else if toFrom < 0 {
		return toFrom + maxSize
	} else {
		return toFrom
	}
}

func MinutesDistance(from, to int) int {

	return unidirectionalRingDistance(from, to, 60)

}

func HoursDistance(from, to int) int {

	return unidirectionalRingDistance(from, to, 24)

}

func FormatTime(hours int, minutes int) string {
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
