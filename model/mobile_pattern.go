package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type MobilePattern struct {
	NID            int       `json:"nid,omitempty" gorm:"column:nID"`
	Pattern        string    `json:"pattern" gorm:"column:pattern"`     // 靓号正则表达式
	TypeName       string    `json:"type_name" gorm:"column:type_name"` // 靓号名称
	IsSort         int       `json:"is_sort" gorm:"column:is_sort"`
	SliceBegin     int       `json:"slice_begin" gorm:"column:slice_begin"` // 起始位
	SliceEnd       int       `json:"slice_end" gorm:"column:slice_end"`     // 结束位
	SortLen        int       `json:"sort_len" gorm:"column:sort_len"`       //号码长度
	MbLevel        int       `json:"mb_level" gorm:"column:mb_level"`       //黑名单等级，1级，100、200、 300级别
	Priority       int       `json:"priority" gorm:"column:priority"`       // 优先级，值越小越优先
	IsRemoveRepeat int       `json:"is_remove_repeat" gorm:"column:is_remove_repeat"`
	IsSlice        int       `json:"is_slice" gorm:"column:is_slice"`
	Hit            int       `json:"hit" gorm:"column:hit"`
	IStatus        int       `json:"i_status" gorm:"column:i_status"` // 状态，0停用，1启动，-1 删除
	JoinDt         time.Time `json:"join_dt" gorm:"column:join_dt"`
	Remark         string    `json:"remark" gorm:"remark"` // 备注
}

func (MobilePattern) TableName() string {
	return "mobile_pattern"
}

type MobilePatternImpl struct {
	DB *gorm.DB
}

var mobilePatternImpl *MobilePatternImpl

type MobilePatternRepo interface {
	GetListByMbLevel(ctx context.Context, mbLevel int) ([]*MobilePattern, error)
}

func GetMobilePatternImpl() MobilePatternRepo {
	return mobilePatternImpl
}

func (m *MobilePatternImpl) GetListByMbLevel(ctx context.Context, mbLevel int) ([]*MobilePattern, error) {
	var res []*MobilePattern
	err := m.DB.WithContext(ctx).Where("mb_level <= ?", mbLevel).Find(&res).Error
	return res, err
}
