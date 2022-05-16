package handlers

import (
	"context"
	"vosBlack/model"
	"vosBlack/proto"
)

func BlackMobileAddHandler(ctx context.Context, req *proto.MobileBlackAddReq, rsp *proto.MobileBlackAddRsp) (err error) {
	blackMobile := &model.MobileBlack{
		Mobile:    req.Mobile,
		MobileAll: req.MobileAll,
		MbLevel:   req.MbLevel,
		GwId:      req.GwId,
		EnID:      req.EnID,
	}
	prefix := req.MobileAll[:3]
	err = model.GetMobileBlackApi().Save(ctx, blackMobile, prefix)
	return err
}

func BlackMobileDelHandler(ctx context.Context, req *proto.BlackMobileDelReq, rsp *proto.BlackMobileDelRsp) (err error) {
	prefix := req.MobileAll[:3]
	err = model.GetMobileBlackApi().Del(ctx, req.NID, prefix)
	return err
}

func BlackMobileListHandler(ctx context.Context, req *proto.BlackMobileListReq, rsp *proto.BlackMobileListRsp) (err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.PageIndex - 1)
	query := model.MobileBlackQueryCondition{
		Limit:  limit,
		Offset: offset,
		Prefix: req.Prefix,
	}
	list, total, err := model.GetMobileBlackApi().GetListByCondition(ctx, query)
	if err != nil {
		return err
	}
	rsp.List = list
	rsp.Total = total
	return nil
}
