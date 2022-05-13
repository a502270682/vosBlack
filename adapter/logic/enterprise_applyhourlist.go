package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"vosBlack/adapter/redis"
)

type ApplyHourListField struct {
	EnID           int    `json:"en_id" gorm:"column:en_id"`
	MbRequestCount int64  `json:"mb_request_count" gorm:"column:mb_request_count"` // 黑名单请求次数
	MbHitCount     int64  `json:"mb_hit_count" gorm:"column:mb_hit_count"`         // 黑名单命中次数
	WnHitCount     int64  `json:"wn_hit_count" gorm:"column:wn_hit_count"`         // 白名单命中次数
	MpRequestCount int64  `json:"mp_request_count" gorm:"column:mp_request_count"` // 靓号请求次数
	MpHitCount     int64  `json:"mp_hit_count" gorm:"column:mp_hit_count"`         // 靓号命中次数
	GwRequestCount int64  `json:"gw_request_count" gorm:"column:gw_request_count"` // 外部接口请求次数
	GwHitCount     int64  `json:"gw_hit_count" gorm:"column:gw_hit_count"`         // 外部接口命中次数
	FqRequestCount int64  `json:"fq_request_count" gorm:"column:fq_request_count"` // 频次请求次数
	FqHitCount     int64  `json:"fq_hit_count" gorm:"column:fq_hit_count"`         // 频次命中次数
	Remark         string `json:"remark" gorm:"column:remark"`                     // 备注
}

func applyHourListCacheKey(enID int) string {
	return fmt.Sprintf("vos_black_enterprise_hour_list:%d", enID)
}

func SetApplyHourListCache(ctx context.Context, enID int, field *ApplyHourListField) error {
	value, err := json.Marshal(field)
	if err != nil {
		return errors.Wrap(err, "marshal apply_field failed")
	}
	err = redis.GetDefaultRedisClient().Set(ctx, applyHourListCacheKey(enID), string(value), 0).Err()
	if err != nil {
		return errors.Wrap(err, "set apply_field failed")
	}
	return nil
}

func GetApplyHourListCache(ctx context.Context, enID int) (*ApplyHourListField, error) {
	result, err := redis.GetDefaultRedisClient().Get(ctx, applyHourListCacheKey(enID)).Result()
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		}
		return nil, err
	}
	var ret *ApplyHourListField
	err = json.Unmarshal([]byte(result), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func DeleteApplyHourListCache(ctx context.Context, enID int) error {
	_, err := redis.GetDefaultRedisClient().Del(ctx, applyHourListCacheKey(enID)).Result()
	if err != nil {
		return err
	}
	return nil
}

// 新建或更新企业近期调用统计
func UpsertEnterpriseApplyHourList(ctx context.Context, enID int, remark string, MbRequestCount, MbHitCount, WnHitCount, MpRequestCount, MpHitCount, GwRequestCount, GwHitCount, FqRequestCount, FqHitCount int64) error {
	field, err := GetApplyHourListCache(ctx, enID)
	if err != nil {
		return errors.Wrap(err, "redis get wrong")
	}
	if field == nil {
		newField := &ApplyHourListField{
			EnID:           enID,
			MbRequestCount: MbRequestCount,
			MbHitCount:     MbHitCount,
			WnHitCount:     WnHitCount,
			MpRequestCount: MpRequestCount,
			MpHitCount:     MpHitCount,
			GwRequestCount: GwRequestCount,
			GwHitCount:     GwHitCount,
			FqRequestCount: FqRequestCount,
			FqHitCount:     FqHitCount,
			Remark:         remark,
		}
		return SetApplyHourListCache(ctx, enID, newField)
	} else {
		field.MbRequestCount += MbRequestCount
		field.MbHitCount += MbHitCount
		field.WnHitCount += WnHitCount
		field.MpRequestCount += MpRequestCount
		field.MpHitCount += MpHitCount
		field.GwRequestCount += GwRequestCount
		field.GwHitCount += GwHitCount
		field.FqRequestCount += FqRequestCount
		field.FqHitCount += FqHitCount
		return SetApplyHourListCache(ctx, enID, field)
	}
}
