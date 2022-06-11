package logic

import (
	"context"
	"testing"
	"vosBlack/utils"
)

//func init() {
//	err := redis.Initialize(&redis.RedisConf{
//		Name:       "default",
//		Addr:       "127.0.0.1:6379",
//		DB:         0,
//		MaxRetries: 3,
//	})
//	if err != nil {
//		panic(err)
//	}
//}

func setCache(ctx context.Context, enID int, today, yesterday, yyesterday string) error {
	err := AddEnterpriseFqCache(ctx, enID, "15201441986", today, 2)
	if err != nil {
		return err
	}
	//err = AddEnterpriseFqCache(ctx, enID, yesterday, 3)
	//if err != nil {
	//	return err
	//}
	err = AddEnterpriseFqCache(ctx, enID, "15201441986", today, 2)
	if err != nil {
		return err
	}
	err = AddEnterpriseFqCache(ctx, enID, "15201441986", today, 2)
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
	phone := "13990163799"

	expireDay := 2

	err := setCache(ctx, enID, today, yesterday, yyesterday)
	if err != nil {
		t.Fatal(err)
	}

	count, err := GetEnterpriseFqFromStartDay(ctx, enID, phone, today, expireDay)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatal("today not correct", count)
	}
	count, err = GetEnterpriseFqFromStartDay(ctx, enID, phone, yesterday, expireDay)
	if err != nil {
		t.Fatal(err)
	}
	if count != 8 {
		t.Fatal("yesterday not correct", count)
	}
	count, err = GetEnterpriseFqFromStartDay(ctx, enID, phone, yyesterday, expireDay)
	if err != nil {
		t.Fatal(err)
	}
	if count != 12 {
		t.Fatal("yyesterday not correct", count)
	}
	t.Log("success")
}
