package common

import "github.com/pkg/errors"

type CtxKey string

var (
	IPCtxKey CtxKey = "ip"
)

const (
	VOSRewrite  = 101
	VOSHttp     = 102
	SVOSHttp    = 201
	DongyunHttp = 301
)

var SignError = errors.New("sign error")
var ReqParamError = errors.New("param error")
var ReqAKError = errors.New("ak error")

const (
	CommonMobileType    = "正常号码"
	SensitiveMobileType = "敏感号码"
	OutMobileType       = "超频号码"
)
