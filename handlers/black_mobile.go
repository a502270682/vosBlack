package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"
	"vosBlack/adapter/error_code"
	"vosBlack/model"
	"vosBlack/proto"

	"gorm.io/gorm"
)

func BlackMobileAddHandler(ctx context.Context, req *proto.MobileBlackAddReq, rsp *proto.MobileBlackAddRsp) *error_code.ReplyError {
	if req.MobileAll == "" || req.Mobile == "" || req.EnID == 0 || req.GwId == 0 {
		return error_code.Error(error_code.CodeParamWrong, fmt.Sprintf("mobile_all(%s) or mobile(%s) or en_id(%d) or gw_id(%d) is needed", req.MobileAll, req.Mobile, req.EnID, req.GwId))
	}
	prefix := req.MobileAll[:3]
	mobile := req.MobileAll[3:]
	if mobile != req.Mobile {
		return error_code.Error(error_code.CodeParamWrong, errors.New("mobile_all and mobile is different").Error())
	}
	mb, err := model.GetMobileBlackApi().GetOneByMobile(ctx, prefix, req.Mobile)
	if err != nil && err != gorm.ErrRecordNotFound {
		return error_code.Error(error_code.CodeSystemError, err.Error())
	}
	if mb != nil {
		// 已存在
		rsp.MobileAll = mb.MobileAll
		return nil
	}
	blackMobile := &model.MobileBlack{
		Mobile:    req.Mobile,
		MobileAll: req.MobileAll,
		MbLevel:   req.MbLevel,
		GwId:      req.GwId,
		EnID:      req.EnID,
		JoinDt:    time.Now(),
	}
	err = model.GetMobileBlackApi().Insert(ctx, blackMobile, prefix)
	if err != nil {
		return error_code.Error(error_code.CodeSystemError, err.Error())
	}
	rsp.MobileAll = blackMobile.MobileAll
	return nil
}

func BlackMobileDelHandler(ctx context.Context, req *proto.BlackMobileDelReq, rsp *proto.BlackMobileDelRsp) *error_code.ReplyError {
	if req.MobileAll == "" {
		return error_code.Error(error_code.CodeParamWrong, "mobile_all is needed")
	}
	prefix := req.MobileAll[:3]
	mobile := req.MobileAll[3:]
	mb, err := model.GetMobileBlackApi().GetOneByMobile(ctx, prefix, mobile)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return error_code.Error(error_code.CodeParamWrong, fmt.Sprintf("black_mobile(%s) not found", req.MobileAll))
		}
		return error_code.Error(error_code.CodeSystemError, err.Error())
	}
	err = model.GetMobileBlackApi().Del(ctx, mb.NID, prefix)
	if err != nil {
		return error_code.Error(error_code.CodeSystemError, err.Error())
	}
	rsp.MobileAll = req.MobileAll
	return nil
}

func BlackMobileListHandler(ctx context.Context, req *proto.BlackMobileListReq, rsp *proto.BlackMobileListRsp) *error_code.ReplyError {
	if req.Prefix == "" {
		return error_code.Error(error_code.CodeParamWrong, "prefix should not be empty")
	}
	if req.PageIndex <= 0 || req.PageSize <= 0 {
		return error_code.Error(error_code.CodeParamWrong, fmt.Sprintf("page_index(%d) or page_size(%d) is invalid", req.PageIndex, req.PageSize))
	}
	if req.PageSize > 1000 {
		req.PageSize = 1000
	}
	limit := req.PageSize
	offset := req.PageSize * (req.PageIndex - 1)
	query := model.MobileBlackQueryCondition{
		Limit:  limit,
		Offset: offset,
		Prefix: req.Prefix,
	}
	list, total, err := model.GetMobileBlackApi().GetListByCondition(ctx, query)
	if err != nil {
		return error_code.Error(error_code.CodeSystemError, err.Error())
	}
	rsp.List = list
	rsp.Total = total
	return nil
}

func BlackMobileInfoHandler(ctx context.Context, req *proto.BlackMobileInfoReq, rsp *proto.BlackMobileInfoRsp) *error_code.ReplyError {
	if req.MobileAll == "" {
		return error_code.Error(error_code.CodeParamWrong, "mobile_all is needed")
	}
	prefix := req.MobileAll[:3]
	mobile := req.MobileAll[3:]
	mb, err := model.GetMobileBlackApi().GetOneByMobile(ctx, prefix, mobile)
	if err != nil || mb == nil {
		return error_code.Error(error_code.CodeSystemError, fmt.Sprintf("black_mobile(%s) not found", req.MobileAll))
	}
	rsp.Res = mb
	return nil
}
