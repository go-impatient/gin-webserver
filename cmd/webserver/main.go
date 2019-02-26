package main

import (
	"fmt"
	"github.com/moocss/go-webserver/src/dao"
	"os"

	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"

	"github.com/moocss/go-webserver/src/app"
	"github.com/moocss/go-webserver/src/config"
	"github.com/moocss/go-webserver/src/pkg/log"
	"github.com/moocss/go-webserver/src/pkg/version"
	"github.com/moocss/go-webserver/src/service"
)

var usageStr = `
              ___.                                           
__  _  __ ____\_ |__   ______ ______________  __ ___________ 
\ \/ \/ // __ \| __ \ /  ___// __ \_  __ \  \/ // __ \_  __ \
 \     /\  ___/| \_\ \\___ \\  ___/|  | \/\   /\  ___/|  | \/
  \/\_/  \___  >___  /____  >\___  >__|    \_/  \___  >__|   
             \/    \/     \/     \/                 \/
`

var flags = []cli.Flag{
	cli.BoolFlag{
		EnvVar: "WEBSERVER_DEBUG",
		Name:   "debug",
		Usage:  "enable app debug mode",
	},
	cli.StringFlag{
		EnvVar: "WEBSERVER_CONFING",
		Name:   "config, c",
		Usage:  "set config file",
	},
}

func start(c *cli.Context) error {
	var g errgroup.Group

	// 设置默认配置
	cfg, err := config.Init(c.String("c"))
	if err != nil {
		log.Infof("Load yaml config file error: '%v'", err)
		os.Exit(-1)
	}

	if c.Bool("debug") {
		cfg.Core.Mode = "dev"
	} else {
		cfg.Core.Mode = "prod"
	}

	// 创建数据库业务相关处理
	dao := dao.New(cfg)

	// 创建封装好的业务服务
	svc := service.New(cfg, dao)

	// 创建后台服务
	app := app.New(cfg, dao, svc)

	// 初始化日志服务
	app.InitLog()

	// 初始化邮件服务
	// app.InitMail()

	g.Go(func() error {
		// 启动服务
		return app.RunHTTPServer()
	})
	g.Go(func() error {
		// 健康检查
		return app.PingServer()
	})

	if err := g.Wait(); err != nil {
		log.Errorf("接口服务停止了：%v", err)
	}

	return g.Wait()
}

func run() {
	app := cli.NewApp()
	app.Name = "webserver"
	app.Version = version.Info.String() // version.Version.String()
	app.Usage = "Golang接口服务器"
	app.UsageText = usageStr
	app.Action = start
	app.Flags = flags
	app.Before = func(c *cli.Context) error {
		fmt.Fprintf(c.App.Writer, "brace for impact\n")
		return nil
	}
	app.After = func(c *cli.Context) error {
		fmt.Fprintf(c.App.Writer, "did we lose anyone?\n")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	run()
}
