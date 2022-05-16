package logic

import (
	"context"
	"testing"
	"vosBlack/adapter/redis"
	"vosBlack/config"
	"vosBlack/utils"
)

func init() {
	config.SetConfigForTest(config.Config{
		FqCountSavedDay: 1,
	})
	err := redis.Initialize(&redis.RedisConf{
		Name:       "default",
		Addr:       "127.0.0.1:6379",
		DB:         0,
		MaxRetries: 3,
	})
	if err != nil {
		panic(err)
	}
}

func setCache(ctx context.Context, enID int, today, yesterday, yyesterday string) error {
	err := AddEnterpriseFqCache(ctx, enID, today, 2)
	if err != nil {
		return err
	}
	//err = AddEnterpriseFqCache(ctx, enID, yesterday, 3)
	//if err != nil {
	//	return err
	//}
	err = AddEnterpriseFqCache(ctx, enID, yesterday, 3)
	if err != nil {
		return err
	}
	err = AddEnterpriseFqCache(ctx, enID, yyesterday, 4)
	if err != nil {
		return err
	}
	return nil
}

func TestFqControl(t *testing.T) {
	ctx := context.Background()
	enID := 123
	today := utils.GetLastNDay0TimeStamp(0)
	yesterday := utils.GetLastNDay0TimeStamp(1)
	yyesterday := utils.GetLastNDay0TimeStamp(2)

	err := setCache(ctx, enID, today, yesterday, yyesterday)
	if err != nil {
		t.Fatal(err)
	}

	count, err := GetEnterpriseFqFromStartDay(ctx, enID, today)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatal("today not correct", count)
	}
	count, err = GetEnterpriseFqFromStartDay(ctx, enID, yesterday)
	if err != nil {
		t.Fatal(err)
	}
	if count != 8 {
		t.Fatal("yesterday not correct", count)
	}
	count, err = GetEnterpriseFqFromStartDay(ctx, enID, yyesterday)
	if err != nil {
		t.Fatal(err)
	}
	if count != 12 {
		t.Fatal("yyesterday not correct", count)
	}
	t.Log("success")
}
