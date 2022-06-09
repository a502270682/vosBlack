package utils

import (
	"fmt"
	"testing"
)

func Test_GetCurHourAndMinute(t *testing.T) {
	hour, minute := GetCurHourAndMinute()
	fmt.Printf("%d %d", hour, minute)
}

func Test_IsReachTime(t *testing.T) {
	beginHour := 9
	endHour := 10
	endMinute := 59
	isReach := IsReachTime(beginHour, 0, endHour, endMinute)
	fmt.Printf("%v", isReach)
}
