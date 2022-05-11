package handlers

import (
	"context"
	"vosBlack/proto"
)

// 外部号码改写规则
func BlackCheckHandler(ctx context.Context, req *proto.BlackCheckReq) (rsp *proto.BlackCheckRsp, err error) {
	//ip := getIP(ctx)

	rsp = &proto.BlackCheckRsp{}
	return rsp, nil
}

//防护接口
func BlackScreeningHandler(ctx context.Context, req *proto.BlackScreeningReq) (rsp *proto.BlackScreeningRsp, err error) {
	rsp = &proto.BlackScreeningRsp{}
	return rsp, nil
}

// 东云接口
func BlackDongYunHandler(ctx context.Context, req *proto.BlackDongYunReq) (rsp *proto.BlackDongYunRsp, err error) {
	rsp = &proto.BlackDongYunRsp{}
	return rsp, nil
}
