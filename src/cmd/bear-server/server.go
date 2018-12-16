package main

import (
	"github.com/moocss/go-webserver/src/config"
	"github.com/moocss/go-webserver/src/log"
	"github.com/moocss/go-webserver/src/server"
	"github.com/moocss/go-webserver/src/storer"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"
)

var flags = []cli.Flag{
	cli.BoolFlag{
		EnvVar: "BEAR_DEBUG",
		Name:   "debug",
		Usage:  "enable server debug mode",
	},
	cli.BoolFlag{
		EnvVar: "BEAR_CONFING",
		Name:   "c",
		Usage:  "set config file",
	},
}

func start(c *cli.Context) error {
	var (
		err error
		g errgroup.Group
	)

	// 初始化数据
	storer.DB.Init()

	// 设置默认参数.
	config.InitConfig(c.String("c"))

	// overwrite server port and address
	if c.String("port") != "" {
		config.Bear.C.Core.Port = c.String("port")
	}
	if c.String("host") != "" {
		config.Bear.C.Core.Host = c.String("host")
	}

	g.Go(func() error {
		// 启动服务
		return server.RunHTTPServer()
	})
	g.Go(func() error {
		// 健康检查
		return server.PingServer()
	})

	if err = g.Wait(); err != nil {
		log.Error("接口服务出错了：", err)
	}

	return g.Wait()
}
