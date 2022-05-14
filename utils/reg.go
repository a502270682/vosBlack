package utils

import "regexp"

var reReqCallee = regexp.MustCompile(`(a(?P<nicelevel>\d+)_)?(b(?P<blacklevel>.*)_)?(?P<realcallee>(?P<prefix>(?U)\d*)(?P<mobile>1[3-9]{1}\d{9})?)$`)

func FindStringSubmatch(str string) []string {
	return reReqCallee.FindStringSubmatch(str)
}
