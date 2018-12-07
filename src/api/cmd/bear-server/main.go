package main

import (
	"fmt"
	"os"

	"github.com/moocss/apiserver/src"
	"github.com/moocss/go-webserver/src/api/bootstrap"
	config "github.com/moocss/go-webserver/src/api/confIg"
	"github.com/moocss/go-webserver/src/api/log"
	"github.com/moocss/go-webserver/src/api/pkg/version"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	rootCmd = &cobra.Command{
		Use:               "bear-server",
		Short:             "bear API server",
		Long:              `Start bear API server`,
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	}
	startCmd = &cobra.Command{
		Use:     "start",
		Short:   "Start bear API server",
		Example: "bear-server start -c src/config/conf.yaml",
		RunE:    start,
	}
)

func init() {
	cobra.OnInitialize(config.Init)

	startCmd.PersistentFlags().StringVarP(&bootstrap.Args.ConfigFile, "conf", "c", "src/conf/conf.yaml",
		"Start server with provided configuration file")
	startCmd.PersistentFlags().StringVarP(&bootstrap.Args.Host, "host", "H", "0.0.0.0",
		"Start server with provided host")
	startCmd.PersistentFlags().IntVarP(&bootstrap.Args.Port, "port", "p", 50000,
		"Start server with provided port")

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(version.Command())
}

func start(_ *cobra.Command, _ []string) error {
	// 解析配置
	parseConfig()

	// 初始日志
	log.Init()

	// 初始数据库

	// 启动服务
	var g errgroup.Group
	g.Go(func() error {
		// 启动服务
		return src.RunHTTPServer()
	})
	g.Go(func() error {
		// 健康检查
		return src.PingServer()
	})

	if err = g.Wait(); err != nil {
		log.Error("接口服务出错了：", err)
	}

	return nil
}

func parseConfig() {

}

func main() {
	// Execute executes the root command.
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
