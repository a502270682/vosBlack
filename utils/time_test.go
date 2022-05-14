package utils

import (
	"fmt"
	"testing"
)

func Test_GetCurHourAndMinute(t *testing.T) {
	hour, minute := GetCurHourAndMinute()
	fmt.Printf("%d %d", hour, minute)
}
