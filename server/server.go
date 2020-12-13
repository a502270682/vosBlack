package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"myGo/adapter/log"
	"myGo/config"
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

func NewServer(ctx context.Context) *Server {
	s := &Server{
		App: cli.NewApp(),
	}
	s.Flags = []cli.Flag{cli.StringFlag{Name: "c", Usage: "Configuration file"}}
	s.Action = func(c *cli.Context) error {
		if c.GlobalString("c") == "" {
			return errors.New("usage: my_go -c configfilepath")
		}
		log.Info(ctx, "start read config: ", c.GlobalString("c"))
		conf, err := initWithConfig(ctx, c.GlobalString("c"))
		if err != nil {
			return errors.Wrap(err, "fail to init conf")
		}
		log.Infof(ctx, "init config success. conf:%+v", conf)
		s.config = conf
		r := gin.Default()
		routes(r)
		if err = r.Run(s.config.HttpPort); err != nil {
			return errors.Wrap(err, "fail to run")
		}

		return nil
	}
	return s
}
