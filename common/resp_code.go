package common

type RespCode int
type RespStatus int
type RequestCommonErrorCode int

const (
	InvalidParam        RequestCommonErrorCode = 400
	MethodNotFound      RequestCommonErrorCode = 404
	InternalServerError RequestCommonErrorCode = 500
)

const (
	StatusOK                 RespStatus = 12000 // 正常呼叫
	BlackMobile              RespStatus = 12001 //黑名单手机号
	NoBalance                RespStatus = 12100 // 欠费/余额不足
	NoAvailableUser          RespStatus = 12002 // 不可呼叫用户
	OutOfFrequency           RespStatus = 12003 // 超频
	IrregularNum             RespStatus = 12004 // 不规则号码
	UnReachTime              RespStatus = 12005 // 禁止呼出时间段
	PrettyNumber             RespStatus = 12006 // 靓号
	DynamicProtect           RespStatus = 12007 // 动态防护
	SystemGatewayBlackMobile RespStatus = 12011 // 第三方黑名单
	NotFound                 RespStatus = -1    // 未找到
	SystemInternalError      RespStatus = 500   // 系统错误
)

var respMsgMap = map[RespStatus]string{
	StatusOK:                 "正常呼叫",
	NoBalance:                "欠费/余额不足",
	NoAvailableUser:          "不可呼叫用户",
	OutOfFrequency:           "超频",
	IrregularNum:             "不规则号码",
	UnReachTime:              "禁止呼出时间段",
	PrettyNumber:             "靓号",
	DynamicProtect:           "动态防护",
	NotFound:                 "未找到",
	SystemInternalError:      "系统错误",
	BlackMobile:              "黑名单手机号",
	SystemGatewayBlackMobile: "第三方黑名单",
}

const (
	Success        = 1
	RespError      = 0
	ParamError     = 4001
	ParamTypeError = 4005
	AKError        = 4000
	CompanyIDError = 4002
	IPError        = 4003
	SignErrorResp  = 4010
	NoBalanceResp  = 4004
	OutOfLimit     = 4011
)

func (r RespCode) Int() int {
	return int(r)
}

func (r RespStatus) Int() int {
	return int(r)
}
