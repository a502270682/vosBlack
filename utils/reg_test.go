package utils

import (
	"fmt"
	"testing"
)

func Test_FindStringSubmatch(t *testing.T) {
	str := "71615121603408"
	matches := FindStringSubmatch(str)
	fmt.Printf("%+v", matches)
}
