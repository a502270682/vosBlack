package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"vosBlack/proto"
)

type SysGatewayDongYun struct {
	AK     string `json:"ak"`
	CallID string `json:"callId"`
	Caller string `json:"caller"`
	Callee string `json:"callee"`
	Sign   string `json:"sign"`
}

func (dy *SysGatewayDongYun) Request(ctx context.Context) (*proto.BlackDongYunRsp, error) {
	uri := "http://bforbid.bsats.cn/bforbid.php"
	res := &proto.BlackDongYunRsp{}
	err := PostJson(ctx, uri, dy, res)
	return res, err
}

type HuaXinBlackCheck struct {
	CallID string `json:"callId"`
	Caller string `json:"caller"`
	Callee string `json:"callee"`
}

type HuaXinResponse struct {
	CallID  string `json:"callId"`
	ForbID  int    `json:"forbID"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (hx *HuaXinBlackCheck) Request(ctx context.Context) (*HuaXinResponse, error) {
	url := "http://api.bjhdsz.cn/api/v2.0/black/check"
	res := &HuaXinResponse{}
	err := PostJson(ctx, url, hx, res)
	return res, err
}

type HuaXinBlackRewrite struct {
	RewriteE164Req struct {
		CallID     string `json:"callId"`     // 通话ID
		CallerE164 string `json:"callerE164"` // 主叫号码
		CalleeE164 string `json:"calleeE164"` // 被叫号码
	} `json:"RewriteE164Req"`
}

func (br *HuaXinBlackRewrite) Request(ctx context.Context) (*proto.BlackCheckRsp, error) {
	url := "http://api.bjhdsz.cn/api/v2.0/black/rewrite"
	res := &proto.BlackCheckRsp{}
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
	err = json.Unmarshal(body, res)
	return err
}
