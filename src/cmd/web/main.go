package main

import (
	"fmt"
	"os"

	"github.com/moocss/go-webserver/src/pkg/version"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "bear-server"
	app.Version = version.Info.String() // version.Version.String()
	app.Usage = "bear server"
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
		os.Exit(-1)
	}
}
