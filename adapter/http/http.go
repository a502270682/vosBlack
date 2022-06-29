package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"vosBlack/adapter/log"
	"vosBlack/proto"
)

type SysGatewayDongYun struct {
	AK     string      `json:"ak"`
	CallID interface{} `json:"callId"`
	Caller string      `json:"caller"`
	Callee string      `json:"callee"`
	Sign   string      `json:"sign"`
}

func (dy *SysGatewayDongYun) Request(ctx context.Context, url string) (*proto.BlackDongYunRsp, error) {
	res := &proto.BlackDongYunRsp{}
	err := PostJson(ctx, url, dy, res)
	return res, err
}

type HuaXinBlackCheck struct {
	CallID interface{} `json:"callId"`
	Caller string      `json:"caller"`
	Callee string      `json:"callee"`
}

type HuaXinResponse struct {
	CallID  interface{} `json:"callId"`
	ForbID  int         `json:"forbID"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}

func (hx *HuaXinBlackCheck) Request(ctx context.Context, url string) (*HuaXinResponse, error) {
	res := &HuaXinResponse{}
	err := PostJson(ctx, url, hx, res)
	return res, err
}

type VOSHttpReq struct {
	CallID interface{} `json:"callId"`
	Caller string      `json:"caller"`
	Callee string      `json:"callee"`
}

type VOSHttpRsp struct {
	CallID  interface{} `json:"callId"`
	ForbID  int         `json:"forbID"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}

func (hx *VOSHttpReq) Request(ctx context.Context, url string) (*VOSHttpRsp, error) {
	res := &VOSHttpRsp{}
	err := PostJson(ctx, url, hx, res)
	return res, err
}

type HuaXinBlackRewrite struct {
	RewriteE164Req struct {
		CallID     interface{} `json:"callId"`     // 通话ID
		CallerE164 string      `json:"callerE164"` // 主叫号码
		CalleeE164 string      `json:"calleeE164"` // 被叫号码
	} `json:"RewriteE164Req"`
}

type HuaXinBlackRewriteRsp struct {
	RewriteE164Rsp struct {
		CallID     interface{} `json:"callId"`
		CallerE164 string      `json:"callerE164"`
		CalleeE164 string      `json:"calleeE164"`
		Code       int         `json:"code"`
		Memo       string      `json:"memo"`
		Status     int         `json:"status"`
	} `json:"RewriteE164Rsp"`
}

func (br *HuaXinBlackRewrite) Request(ctx context.Context, url string) (*HuaXinBlackRewriteRsp, error) {
	res := &HuaXinBlackRewriteRsp{}
	err := PostJson(ctx, url, br, res)
	return res, err
}
func PostJson(ctx context.Context, url string, req, res interface{}) error {
	post, err := json.Marshal(req)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(post))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Infof(ctx, "post json response: %s", string(body))
	err = json.Unmarshal(body, res)
	return err
}
