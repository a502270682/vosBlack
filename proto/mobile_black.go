package proto

import "vosBlack/model"

type MobileBlackAddReq struct {
	Mobile    string `json:"mobile" gorm:"column:mobile"`     // 手机号码，后8位
	MobileAll string `json:"mobile_all" gorm:"column:mobile"` // 完整手机号码
	MbLevel   int    `json:"mb_level" gorm:"column:mb_level"` // 黑名单级别
	GwId      int    `json:"gw_id" gorm:"column:gw_id"`       // 调用网关id
	EnID      int    `json:"en_id" gorm:"column:en_id"`       // 请求企业ID
}

type MobileBlackAddRsp struct {
}

type BlackMobileDelReq struct {
	NID       int    `json:"nid" form:"nid"`
	MobileAll string `json:"mobile_all" form:"mobile_all"`
}

type BlackMobileDelRsp struct {
}

type BlackMobileListReq struct {
	PageIndex int    `json:"page_index,omitempty" form:"page_index"`
	PageSize  int    `json:"page_size,omitempty" form:"page_size"`
	Prefix    string `json:"prefix" form:"prefix"`
}
type BlackMobileListRsp struct {
	List  []*model.MobileBlack `json:"list"`
	Total int64                `json:"total"`
}
