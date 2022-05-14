package model

import "time"

type SysGateway struct {
	NID         int       `json:"nid,omitempty" gorm:"column:nID"`
	GwName      string    `json:"gw_name" gorm:"column:gw_name"`           // 网关名称
	GwUrl       string    `json:"gw_url" gorm:"column:gw_url"`             //网关调用地址
	GwType      int       `json:"gw_type" gorm:"column:gw_type"`           // 网关类型
	Priority    int       `json:"priority" gorm:"column:priority"`         // 优先级
	BlackPrefix string    `json:"black_prefix" gorm:"column:black_prefix"` // 接口其他参数
	IStatus     int       `json:"i_status" gorm:"i_status"`                //状态，0停用，1启动，-1删除
	MbLevel     int       `json:"mb_level" gorm:"mb_level"`                // 黑名单等级
	JoinDt      time.Time `json:"join_dt" gorm:"join_dt"`                  // 添加时间
	Remark      string    `json:"remark" gorm:"remark"`                    // 备注
}

func (SysGateway) TableName() string {
	return "sys_gateway"
}
