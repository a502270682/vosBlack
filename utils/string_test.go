package utils

import (
	"fmt"
	"testing"
)

func Test_GetPhone(t *testing.T) {
	callee := "8021213810507903"
	prefix, realCallee, phoneType := GetPhone(callee)
	fmt.Printf("prefix: %s realCallee : %s phoneType : %d \n", prefix, realCallee, phoneType)
}
