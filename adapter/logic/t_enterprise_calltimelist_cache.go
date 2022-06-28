package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"vosBlack/adapter/log"
	"vosBlack/adapter/redis"
	"vosBlack/model"
)

func enterpriseCallTimeListKey(ipID int, blackID int) string {
	return fmt.Sprintf(vosBlackEnterpriseCallTimeListKey, ipID, blackID)
}

func GetEnterpriseCallTimeListWithCache(ctx context.Context, enID int, blackID int) ([]*model.EnterpriseCalltimelist, error) {
	res, err := redis.GetDefaultRedisClient().Get(ctx, enterpriseCallTimeListKey(enID, blackID)).Result()
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	if err == redis.ErrNil {
		// redis没有，从db里取
		calltimelist, err := model.GetEnterpriseCalltimelistImpl().GetByEnID(ctx, enID, blackID)
		if err != nil {
			return nil, err
		}
		// 如果取到了，存一下缓存
		err = SetEnterpriseCallTimeListCache(ctx, enID, blackID, calltimelist)
		if err != nil {
			log.Warnf(ctx, "SetEnterpriseCallTimeListCache failed, err:%s, enID:%d, blackID:%s", err.Error(), enID, blackID)
		}
		return calltimelist, nil
	}
	var ret []*model.EnterpriseCalltimelist
	err = json.Unmarshal([]byte(res), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func SetEnterpriseCallTimeListCache(ctx context.Context, enID int, blackID int, calltimelist []*model.EnterpriseCalltimelist) error {
	res, err := json.Marshal(calltimelist)
	if err != nil {
		return err
	}
	return redis.GetDefaultRedisClient().Set(ctx, enterpriseCallTimeListKey(enID, blackID), string(res), 0).Err()
}
