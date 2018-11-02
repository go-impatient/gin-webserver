package cmd

import (
	"github.com/spf13/cobra"
	"github.com/moocss/go-webserver/src/api/conf"
	"github.com/moocss/go-webserver/src/api/bootstrap"
)

var (
	rootCmd = &cobra.Command{
		Use:               "web-api",
		Short:             "api server",
		Long:              `Start api server`,
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	}
	startCmd = &cobra.Command{
		Use:     "start",
		Short:   "Start api server",
		Example: "web-api start -c src/conf/conf.yaml",
		RunE:    start,
	}
)

func init()  {
	cobra.OnInitialize(conf.Init)

	startCmd.PersistentFlags().StringVarP(&bootstrap.Args.ConfigFile, "conf", "c", "src/conf/conf.yaml",
		"Start server with provided configuration file")
	startCmd.PersistentFlags().StringVarP(&bootstrap.Args.Host, "host", "H", "0.0.0.0",
		"Start server with provided host")
	startCmd.PersistentFlags().IntVarP(&bootstrap.Args.Port, "port", "p", 50000,
		"Start server with provided port")

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(versionCmd())
}

func start(_ *cobra.Command, _ []string) error {

	return nil
}

// Execute executes the root command.
func Execute() {
	rootCmd.Execute()
}
