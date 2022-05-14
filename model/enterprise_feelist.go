package model

import "time"

// EnterpriseFeeList 企业账号费用情况
type EnterpriseFeeList struct {
	NID       int       `json:"nid,omitempty" gorm:"column:nID"`
	EnID      int       `json:"en_id,omitempty" gorm:"column:en_id"`
	FeeRate   float64   `json:"fee_rate,omitempty" gorm:"column:fee_rate"`
	FeeIncome float64   `json:"fee_income,omitempty" gorm:"column:fee_income"`
	FeePayout float64   `json:"fee_payout,omitempty" gorm:"column:fee_payout"`
	FeeCredit int       `json:"fee_credit,omitempty" gorm:"column:fee_credit"`
	JoinDt    time.Time `json:"join_dt" gorm:"column:join_dt"`
}

func (EnterpriseFeeList) TableName() string {
	return "t_enterprise_feelist"
}
