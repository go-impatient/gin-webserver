package main

import (
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"
)

func server(c *cli.Context) error {

	var g errgroup.Group

	g.Go(func() error {
		// 启动服务

		return nil
	})
	g.Go(func() error {
		// 健康检查

		return nil
	})

	return g.Wait()
}
