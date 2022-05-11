package proto

type BlackScreeningReq struct {
	CallID string `json:"callId"` // 呼叫标识
	Callee string `json:"callee"` // 被叫号码
	Caller string `json:"caller"` // 主叫号码
}

type BlackScreeningRsp struct {
	CallID string `json:"callId"`
	ForbID int    `json:"forbid"`
	Status int    `json:"status"`
}

type BlackCheckReq struct {
	RewriteE164Req struct {
		CallID     int    `json:"callId"`     // 通话ID
		CallerE164 string `json:"callerE164"` // 主叫号码
		CalleeE164 string `json:"calleeE164"` // 被叫号码
	} `json:"RewriteE164Req"`
}

type BlackCheckRsp struct {
	RewriteE164Rsp struct {
		Status     int    `json:"status"`
		CallID     int    `json:"callId"`
		CallerE164 string `json:"callerE164"`
		CalleeE164 string `json:"calleeE164"`
	} `json:"RewriteE164Rsp"`
}

type BlackDongYunReq struct {
	AK     string `json:"ak"`     // 必填，接口提供方提供的【企业id】
	CallID string `json:"callId"` // 呼叫唯一标识， 必填，必须唯一
	Caller string `json:"caller"` // 主叫号码， 必填，不得低于3位
	Callee string `json:"callee"` // 被叫号码， 必填，用半角逗号(,)隔开的号码清单
	Sign   string `json:"sign"`   // 签名， 必填, 生成算法
}

type BlackDongYunRsp struct {
	Code          int                  `json:"code"`          // 响应状态
	Msg           string               `json:"msg"`           // 响应描述,success:成功 ,fail:失败
	CallID        string               `json:"callId"`        // 通话ID
	TransactionID string               `json:"transactionId"` // 主叫号码
	List          []BlackDongYunDetail `json:"list"`          // 响应对象数组
	Status        int                  `json:"status"`        //
}

type BlackDongYunDetail struct {
	Mobile string `json:"mobile"` // 被叫号码
	Forbid string `json:"forbid"` // 0是正常号码 1是敏感号码 2是超频号码
	Msg    string `json:"msg"`
}
