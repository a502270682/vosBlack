package logic

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"vosBlack/adapter/log"
	"vosBlack/adapter/redis"
	"vosBlack/utils"

	"github.com/pkg/errors"
)

//type EnterpriseFrequencyInfo struct {
//	FqRequestCount int64
//}

// key:enterprise_id field:day(ex:2006-01-02) value:5(已访问次数)
func enterpriseFqHashMapKey(enID int, phone string) string {
	return fmt.Sprintf(vosBlackEnterpriseFrequencyKey, enID, phone)
}

func syncHashMapKey(enID int) string {
	return fmt.Sprintf(syncVosBlackEnterpriseFrequencyKey, enID)
}

func GetEnterpriseFqCache(ctx context.Context, enID int, phone string, dayStamp string) (int64, error) {
	count, err := redis.GetDefaultRedisClient().HGet(ctx, enterpriseFqHashMapKey(enID, phone), dayStamp).Int64()
	if err != nil {
		if err == redis.ErrNil {
			return 0, nil
		}
		return 0, errors.Wrap(err, fmt.Sprintf("enID(%d) get enterprise_fq_count failed", enID))
	}
	return count, nil
}

// 设置目标日期时间戳该企业请求次数
func AddEnterpriseFqCache(ctx context.Context, enID int, phone string, dayStamp string, count int64) error {
	res, err := GetEnterpriseFqCache(ctx, enID, phone, dayStamp)
	if err != nil {
		return err
	}
	err = redis.GetDefaultRedisClient().HSet(ctx, enterpriseFqHashMapKey(enID, phone), dayStamp, res+count).Err()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("enID(%d) set enterprise_fq_count failed", enID))
	}
	redis.GetDefaultRedisClient().Expire(ctx, enterpriseFqHashMapKey(enID, phone), 86400*7*time.Second)
	return nil
}

// 获取从目标日期时间戳开始记录的该企业请求次数
func GetEnterpriseFqFromStartDay(ctx context.Context, enID int, phone string, startDayStamp string, callCycle int) (int64, error) {
	redis.Lock(syncHashMapKey(enID))
	defer redis.Unlock(syncHashMapKey(enID))
	m, err := redis.GetDefaultRedisClient().HGetAll(ctx, enterpriseFqHashMapKey(enID, phone)).Result()
	if err != nil {
		if err == redis.ErrNil {
			return 0, nil
		}
		return 0, err
	}
	total := int64(0)
	expireDay := utils.GetLastNDay0TimeStamp(callCycle)
	for day, count := range m {
		c, _ := strconv.ParseInt(count, 10, 64)
		if day >= startDayStamp {
			total += c
		} else if expireDay > day {
			// delete no use day
			err = DeleteOutOfDateEnterpriseFqField(ctx, enID, phone, day)
			if err != nil {
				log.Warnf(ctx, fmt.Sprintf("DeleteOutOfDateEnterpriseFqField enID(%d)_day(%s) err:%+v", enID, day, err))
				continue
			}
		}
	}
	return total, nil
}

func DeleteOutOfDateEnterpriseFqField(ctx context.Context, enID int, phone string, dayStamp string) error {
	return redis.GetDefaultRedisClient().HDel(ctx, enterpriseFqHashMapKey(enID, phone), dayStamp).Err()
}
