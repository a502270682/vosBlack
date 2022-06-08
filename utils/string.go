package utils

import (
	"sort"
	"strings"
)

func StringContains(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	//index的取值：[0,len(str_array)]
	if index < len(str_array) && str_array[index] == target { //需要注意此处的判断，先判断 &&左侧的条件，如果不满足则结束此处判断，不会再进行右侧的判断
		return true
	}
	return false
}

var SP = []string{"013", "014", "015", "016", "017", "018", "019"}
var SP1 = []string{"13", "14", "15", "16", "17", "18", "19"}

func GetPhone(callee string) (prefix, realCallee string, phoneType int) {
	if strings.HasPrefix(callee, "0") {
		if StringContains(callee[:3], SP) && len(callee) == 12 {
			return "ALL", callee[1:], 1
		}
		return "ALL", callee, 0
	}
	if strings.HasPrefix(callee, "1") {
		if StringContains(callee[:2], SP1) && len(callee) == 11 {
			return "ALL", callee, 1
		}
		return "ALL", callee, 0

	}
	matchs := FindStringSubmatch(callee)
	if matchs[6] != "" && matchs[7] != "" {
		return matchs[6], matchs[7], 1
	}
	prefix = callee[:3]
	realCallee = callee[3:]
	return prefix, realCallee, 0
}
