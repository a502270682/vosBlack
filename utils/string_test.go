package utils

import (
	"fmt"
	"testing"
)

func Test_GetPhone(t *testing.T) {
	callee := "80615201441986"
	prefix, realCallee, phoneType := GetPhone(callee)
	fmt.Printf("prefix: %s realCallee : %s phoneType : %d \n", prefix, realCallee, phoneType)
}
