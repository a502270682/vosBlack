package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"vosBlack/adapter/log"
	"vosBlack/adapter/redis"
	"vosBlack/model"
)

func mobilePatternKey(mbLevel int) string {
	return fmt.Sprintf(vosBlackMobilePatternKey, mbLevel)
}

func GetMobilePatternWithCache(ctx context.Context, mbLevel int) ([]*model.MobilePattern, error) {
	res, err := redis.GetDefaultRedisClient().Get(ctx, mobilePatternKey(mbLevel)).Result()
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	if err == redis.ErrNil {
		// redis没有，从db里取
		mpList, err := model.GetMobilePatternImpl().GetListByMbLevel(ctx, mbLevel)
		if err != nil {
			return nil, err
		}
		// 如果取到了，存一下缓存
		err = SetMobilePatternCache(ctx, mbLevel, mpList)
		if err != nil {
			log.Warnf(ctx, "SetMobilePatternCache failed, err:%s, mbLevel:%d", err.Error(), mbLevel)
		}
		return mpList, nil
	}
	var ret []*model.MobilePattern
	err = json.Unmarshal([]byte(res), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func SetMobilePatternCache(ctx context.Context, mbLevel int, mpList []*model.MobilePattern) error {
	res, err := json.Marshal(mpList)
	if err != nil {
		return err
	}
	return redis.GetDefaultRedisClient().Set(ctx, mobilePatternKey(mbLevel), string(res), 0).Err()
}
