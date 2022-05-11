package handlers

import (
	"context"
	"vosBlack/adapter/error_code"
	"vosBlack/proto"
)

func PingHandler(ctx context.Context, req *proto.PingReq, rsp *proto.PingRsp) *error_code.ReplyError {
	rsp.Success = "hello"
	return nil
}
