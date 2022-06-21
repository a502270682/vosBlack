package utils

import (
	"fmt"
	"testing"
	"time"
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

func Test_GetHour(t *testing.T) {
	str := "2022-06-21 00:00:00"
	loc, _ := time.LoadLocation("Asia/Shanghai")
	dateTime, _ := time.ParseInLocation("2006-01-02 15:04:05", str, loc)
	hour := dateTime.Hour()
	if dateTime.Hour() == 0 {
		hour = 23
	} else {
		hour = hour - 1
	}
	fmt.Printf("%d", hour)
}
