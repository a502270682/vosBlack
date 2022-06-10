package server

import (
	"context"
	"fmt"
	"vosBlack/adapter/log"
	"vosBlack/adapter/mysql"
	"vosBlack/adapter/redis"
	"vosBlack/config"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type Server struct {
	*cli.App
	config *config.Config
}

func initWithConfig(ctx context.Context, filePath string) (*config.Config, error) {
	conf, err := config.Load(filePath)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func initLog(ctx context.Context, conf *config.Config) {
	log.NewLoggerWithOptions(conf.Log)
}

func initMysql(conf *config.Config) error {
	db, err := mysql.InitializeMainDb(conf.Mysql.Master)
	if err != nil {
		return err
	}
	mysql.InitEntityDao(db)
	return nil
}

func NewServer(ctx context.Context) *Server {
	s := &Server{
		App: cli.NewApp(),
	}
	s.Flags = []cli.Flag{cli.StringFlag{Name: "c", Usage: "Configuration file"}}
	s.Action = func(c *cli.Context) error {
		if c.GlobalString("c") == "" {
			return errors.New("usage: vos_black -c configfilepath")
		}

		fmt.Printf("start read config: %v", c.GlobalString("c"))
		conf, err := initWithConfig(ctx, c.GlobalString("c"))
		if err != nil {
			return errors.Wrap(err, "fail to init conf")
		}
		fmt.Printf("init config success. conf:%+v", conf)
		s.config = conf
		initLog(context.Background(), s.config)
		// mysql
		err = initMysql(conf)
		if err != nil {
			return errors.Wrap(err, "fail to init mysql")
		}

		// redis init
		err = redis.Initialize(conf.Redis.Default)
		if err != nil {
			return err
		}

		r := gin.Default()
		routes(r)
		if err = r.Run(s.config.HTTPPort); err != nil {
			return errors.Wrap(err, "fail to run")
		}

		return nil
	}
	return s
}
