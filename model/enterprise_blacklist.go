package model

// EnterpriseBlacklist 企业的黑名单规则设置表
type EnterpriseBlacklist struct {
	NID           int    `json:"nid,omitempty" gorm:"column:nID"`
	EnID          int    `json:"en_id,omitempty" gorm:"column:en_id"`         // 请求企业ID
	EnIPID        int    `json:"en_ip_id" gorm:"column:en_ip_id"`             // 请求企业的IP的ID，0.0.0.0的代表全部
	Qianzhui      string `json:"qianzhui" gorm:"column:qianzhui"`             // 号码前缀，ALL代表全部，先完全匹配，匹配不上的用ALL的设置
	PatternLevel  int    `json:"pattern_level" gorm:"column:pattern_level"`   // 靓号限制[100到400]，-1代表不启用。
	IsWhitenum    int    `json:"is_whitenum" gorm:"column:is_whitenum"`       // 是否启用白名单功能
	IsCalltime    int    `json:"is_calltime" gorm:"column:is_calltime"`       // 是否启用呼叫时间段限制
	BlacknumLevel int    `json:"blacknum_level" gorm:"column:blacknum_level"` // 本地黑名单等级，-1代表不启用。
	GatewayLevel  int    `json:"gateway_level" gorm:"column:gateway_level"`   // 调用第三方网关的网关编号，-1代表不启用
	IsFrequency   int    `json:"is_frequency" gorm:"column:is_frequency"`     // 是否启用频次限制
	CallCycle     int    `json:"call_cycle" gorm:"column:call_cycle"`         // 呼叫周期，多少个小时，大于0，-1代表不启用
	CallCount     int    `json:"call_count" gorm:"column:call_count"`         // 周期内，可以呼叫的次数，大于0，-1不限制次数
	ConnCycle     int    `json:"conn_cycle" gorm:"column:conn_cycle"`         // 接通周期，多少个小时，大于0，-1代表不启用
	ConnCount     int    `json:"conn_count" gorm:"column:conn_count"`         // 周期内，可以接通的次数，大于0，-1不限制次数
}

func (EnterpriseBlacklist) TableName() string {
	return "t_enterprise_blacklist"
}
