package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"vosBlack/adapter/log"
	"vosBlack/adapter/redis"
	"vosBlack/model"
)

func mobileWhiteNumKey(enID int, whiteNum string) string {
	return fmt.Sprintf(vosBlackMobileWhitenumKey, enID, whiteNum)
}

func GetMobileWhiteNumWithCache(ctx context.Context, enID int, whiteNum string) (*model.MobileWhitenum, error) {
	res, err := redis.GetDefaultRedisClient().Get(ctx, mobileWhiteNumKey(enID, whiteNum)).Result()
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	if err == redis.ErrNil {
		// redis没有，从db里取
		whitenum, err := model.GetMobileWhitenumImpl().GetByWhiteNum(ctx, enID, whiteNum)
		if err != nil {
			return nil, err
		}
		// 如果取到了，存一下缓存
		err = SetMobileWhiteNumCache(ctx, enID, whiteNum, whitenum)
		if err != nil {
			log.Warnf(ctx, "SetMobileWhiteNumCache failed, err:%s, enID:%d, whiteNum:%s", err.Error(), enID, whiteNum)
		}
		return whitenum, nil
	}
	var ret *model.MobileWhitenum
	err = json.Unmarshal([]byte(res), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func SetMobileWhiteNumCache(ctx context.Context, enID int, whiteNum string, whitenum *model.MobileWhitenum) error {
	res, err := json.Marshal(whitenum)
	if err != nil {
		return err
	}
	return redis.GetDefaultRedisClient().Set(ctx, mobileWhiteNumKey(enID, whiteNum), string(res), 0).Err()
}
