package service

import (
	"vosBlack/adapter/mysql"
	"vosBlack/adapter/redis"
	"vosBlack/config"
)

func init() {
	err := redis.Initialize(&redis.RedisConf{
		Name:       "default",
		Addr:       "127.0.0.1:6379",
		DB:         0,
		MaxRetries: 3,
	})
	db, err := mysql.InitializeMainDb(config.ConnectionConfig{
		User:     "root",
		Password: "LUbin123!",
		Host:     "bj-cynosdbmysql-grp-pofd9c6u.sql.tencentcdb.com",
		Port:     22551,
		Db:       "test",
		MaxIdle:  5,
		MaxOpen:  10,
		Debug:    true,
	})
	if err != nil {
		panic(err)
	}
	mysql.InitEntityDao(db)
	if err != nil {
		panic(err)
	}
}
