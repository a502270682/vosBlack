package main

import (
	"context"
	"flag"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vosBlack/adapter/log"
	"vosBlack/adapter/mysql"
	"vosBlack/adapter/redis"
	"vosBlack/config"
	"vosBlack/cron"
)

type Server struct {
	*cli.App
	config *config.Config
}

var flagConfigPath = flag.String("config", "./vos_black.toml", "")

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()

	log.Info(ctx, "start read config: ", *flagConfigPath)
	conf, err := config.Load(*flagConfigPath)
	if err != nil {
		panic(errors.Wrap(err, "fail to init conf"))
	}
	log.Infof(ctx, "init config success. conf:%+v", conf)
	// mysql
	db, err := mysql.InitializeMainDb(conf.Mysql.Master)
	if err != nil {
		panic(err)
	}
	mysql.InitEntityDao(db)

	// redis init
	err = redis.Initialize(conf.Redis.Default)
	if err != nil {
		panic(err)
	}

	cron.StartCron()
	signals := make(chan os.Signal, 1)
	defer close(signals)
	signal.Notify(signals, os.Kill, os.Interrupt, syscall.SIGBUS, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)
	<-signals
	cron.StopCron()
	log.Info(ctx, "service start to shut down")
}
