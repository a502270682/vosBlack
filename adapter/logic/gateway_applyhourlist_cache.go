package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"vosBlack/adapter/redis"
)

type GwApplyHourListField struct {
	GwID           int    `json:"gw_id" gorm:"column:gw_id"`
	MbRequestCount int64  `json:"mb_request_count" gorm:"column:mb_request_count"` // 黑名单请求次数
	MbHitCount     int64  `json:"mb_hit_count" gorm:"column:mb_hit_count"`         // 黑名单命中次数
	Remark         string `json:"remark" gorm:"column:remark"`                     // 备注
}

func gwApplyHourListCacheKey(gwID int) string {
	return fmt.Sprintf(vosBlackGatewayHourListKey, gwID)
}

func SetGwApplyHourListCache(ctx context.Context, gwID int, field *GwApplyHourListField) error {
	value, err := json.Marshal(field)
	if err != nil {
		return errors.Wrap(err, "marshal gw_apply_field failed")
	}
	err = redis.GetDefaultRedisClient().Set(ctx, gwApplyHourListCacheKey(gwID), string(value), 0).Err()
	if err != nil {
		return errors.Wrap(err, "set gw_apply_field failed")
	}
	return nil
}

func GetGwApplyHourListCache(ctx context.Context, gwID int) (*GwApplyHourListField, error) {
	result, err := redis.GetDefaultRedisClient().Get(ctx, gwApplyHourListCacheKey(gwID)).Result()
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		}
		return nil, err
	}
	var ret *GwApplyHourListField
	err = json.Unmarshal([]byte(result), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func DeleteGwApplyHourListCache(ctx context.Context, gwID int) error {
	_, err := redis.GetDefaultRedisClient().Del(ctx, gwApplyHourListCacheKey(gwID)).Result()
	if err != nil {
		return err
	}
	return nil
}

// 新建或更新网关近期调用统计
func UpsertGwEnterpriseApplyHourList(ctx context.Context, gwID int, remark string, MbRequestCount, MbHitCount int64) error {
	field, err := GetGwApplyHourListCache(ctx, gwID)
	if err != nil {
		return errors.Wrap(err, "redis get wrong")
	}
	if field == nil {
		newField := &GwApplyHourListField{
			GwID:           gwID,
			MbRequestCount: MbRequestCount,
			MbHitCount:     MbHitCount,
			Remark:         remark,
		}
		return SetGwApplyHourListCache(ctx, gwID, newField)
	} else {
		field.MbRequestCount += MbRequestCount
		field.MbHitCount += MbHitCount
		return SetGwApplyHourListCache(ctx, gwID, field)
	}
}
