package handlers

import (
	"context"
	"vosBlack/proto"
)

func BlackMobileAddHandler(ctx context.Context, req *proto.MobileBlackAddReq) (rsp *proto.MobileBlackAddRsp, err error) {
	rsp = &proto.MobileBlackAddRsp{}
	return rsp, nil
}

func BlackMobileDelHandler(ctx context.Context, req *proto.BlackMobileDelReq) (rsp *proto.BlackMobileDelRsp, err error) {
	rsp = &proto.BlackMobileDelRsp{}
	return rsp, nil
}

func BlackMobileListHandler(ctx context.Context, req *proto.BlackMobileListReq) (rsp *proto.BlackMobileListRsp, err error) {
	rsp = &proto.BlackMobileListRsp{}
	return rsp, nil
}
