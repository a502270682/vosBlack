package handlers

import (
	"context"
	"vosBlack/model"
	"vosBlack/proto"
)

func BlackMobileAddHandler(ctx context.Context, req *proto.MobileBlackAddReq) (rsp *proto.MobileBlackAddRsp, err error) {
	rsp = &proto.MobileBlackAddRsp{}
	blackMobile := &model.MobileBlack{
		Mobile:    req.Mobile,
		MobileAll: req.MobileAll,
		MbLevel:   req.MbLevel,
		GwId:      req.GwId,
		EnID:      req.EnID,
	}
	prefix := req.MobileAll[:3]
	err = model.GetMobileBlackApi().Save(ctx, blackMobile, prefix)
	return rsp, err
}

func BlackMobileDelHandler(ctx context.Context, req *proto.BlackMobileDelReq) (rsp *proto.BlackMobileDelRsp, err error) {
	rsp = &proto.BlackMobileDelRsp{}
	prefix := req.MobileAll[:3]
	err = model.GetMobileBlackApi().Del(ctx, req.NID, prefix)
	return rsp, err
}

func BlackMobileListHandler(ctx context.Context, req *proto.BlackMobileListReq) (rsp *proto.BlackMobileListRsp, err error) {
	rsp = &proto.BlackMobileListRsp{}
	limit := req.PageSize
	offset := req.PageSize * (req.PageIndex - 1)
	query := model.MobileBlackQueryCondition{
		Limit:  limit,
		Offset: offset,
		Prefix: req.Prefix,
	}
	list, total, err := model.GetMobileBlackApi().GetListByCondition(ctx, query)
	if err != nil {
		return nil, err
	}
	rsp.List = list
	rsp.Total = total
	return rsp, nil
}
