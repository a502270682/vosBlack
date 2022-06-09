package utils

import (
	"fmt"
	"strconv"
	"time"
)

func GetLastNDay0TimeStamp(n int) string {
	t := time.Now().Add(-time.Hour * 24 * time.Duration(n))
	return strconv.FormatInt(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Unix(), 10)
}

/*
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", "2016-06-13 15:34:39", time.Local)
	// 整点（向下取整）
	fmt.Println(t.Truncate(1 * time.Hour))
	// 整点（最接近）
	fmt.Println(t.Round(1 * time.Hour))

	// 整分（向下取整）
	fmt.Println(t.Truncate(1 * time.Minute))
	// 整分（最接近）
	fmt.Println(t.Round(1 * time.Minute))
	t2, _ := time.ParseInLocation("2006-01-02 15:04:05", t.Format("2006-01-02 15:00:00"), time.Local)
	fmt.Println(t2)
*/

func GetCurHourAndMinute() (hour, minute int) {
	t := time.Now()
	return t.Truncate(1 * time.Hour).Hour(), t.Truncate(1 * time.Minute).Minute()

}

func IsReachTime(beginHour, beginMinute, endHour, endMinute int) bool {
	curDay := time.Now().Format("2006-01-02")
	beginTimeStr := fmt.Sprintf("%s %02d:%02d:00", curDay, beginHour, beginMinute)
	endTimeStr := fmt.Sprintf("%s %02d:%02d:59", curDay, endHour, endMinute)
	var timeLayoutStr = "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Shanghai")
	beginTime, _ := time.ParseInLocation(timeLayoutStr, beginTimeStr, loc)
	endTime, _ := time.ParseInLocation(timeLayoutStr, endTimeStr, loc)
	currentTime := time.Now().In(loc).Unix()
	beginTimeStamp := beginTime.Unix()
	endTimeStamp := endTime.Unix()
	if beginTimeStamp <= currentTime && currentTime <= endTimeStamp {
		return true
	}
	return false
}
