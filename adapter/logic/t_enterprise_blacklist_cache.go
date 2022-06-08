package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"vosBlack/adapter/log"
	"vosBlack/adapter/redis"
	"vosBlack/model"

	"gorm.io/gorm"
)

func enterpriseBlackListKey(ipID int, prefix string) string {
	return fmt.Sprintf(vosBlackEnterpriseBlackListKey, ipID, prefix)
}

func GetEnterpriseBlackListWithCache(ctx context.Context, ipID int, prefix string) (*model.EnterpriseBlacklist, error) {
	res, err := redis.GetDefaultRedisClient().Get(ctx, enterpriseBlackListKey(ipID, prefix)).Result()
	if err != nil {
		if err == redis.ErrNil {
			res, err = redis.GetDefaultRedisClient().Get(ctx, enterpriseBlackListKey(ipID, "ALL")).Result()
			if err != nil && err != redis.ErrNil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	if err == redis.ErrNil {
		// redis没有，从db里取
		blackRule, err := model.GetEnterpriseBlacklistImpl().GetEnterpriseBlacklistByIPAndQianzhui(ctx, ipID, prefix)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if err == gorm.ErrRecordNotFound {
			blackRule, err = model.GetEnterpriseBlacklistImpl().GetEnterpriseBlacklistByIPAndQianzhui(ctx, ipID, "ALL")
			if err != nil {
				return nil, err
			}
			err = SetEnterpriseBlackListCache(ctx, ipID, "ALL", blackRule)
			if err != nil {
				log.Warnf(ctx, "SetEnterpriseBlackListCache failed, err:%s, ipID:%d, prefix:%s", err.Error(), ipID, prefix)
			}
			return blackRule, nil
		}
		// 如果取到了，存一下缓存
		err = SetEnterpriseBlackListCache(ctx, ipID, prefix, blackRule)
		if err != nil {
			log.Warnf(ctx, "SetEnterpriseBlackListCache failed, err:%s, ipID:%d, prefix:%s", err.Error(), ipID, prefix)
		}
		return blackRule, nil
	}
	var ret *model.EnterpriseBlacklist
	err = json.Unmarshal([]byte(res), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func SetEnterpriseBlackListCache(ctx context.Context, ipID int, prefix string, blacklist *model.EnterpriseBlacklist) error {
	res, err := json.Marshal(blacklist)
	if err != nil {
		return err
	}
	return redis.GetDefaultRedisClient().Set(ctx, enterpriseBlackListKey(ipID, prefix), string(res), 0).Err()
}
