package utils

import (
	"fmt"
	"regexp"
	"testing"
)

func Test_FindStringSubmatch(t *testing.T) {
	str := "71615121603408"
	matches := FindStringSubmatch(str)
	fmt.Printf("%+v", matches)
}

func Test_Regxp(t *testing.T) {
	reg := regexp.MustCompile("^\\d+(\\d)\\1{2}(?!\\1)\\d$")
	realCallee := "15121601112"
	if reg.Match([]byte(realCallee)) {
		fmt.Println("匹配成功")
	} else {
		fmt.Println("匹配失败")
	}
}
