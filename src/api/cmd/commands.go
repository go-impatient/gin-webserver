package cmd

import (
	"github.com/spf13/cobra"
	"github.com/moocss/go-webserver/src/api/config"
	"github.com/moocss/go-webserver/src/api/bootstrap"
)

var (
	rootCmd = &cobra.Command{
		Use:               "web-api",
		Short:             "api server",
		Long:              `
							.__                                        
_____  ______ |__| ______ ______________  __ ___________ 
\__  \ \____ \|  |/  ___// __ \_  __ \  \/ // __ \_  __ \
 / __ \|  |_> >  |\___ \\  ___/|  | \/\   /\  ___/|  | \/
(____  /   __/|__/____  >\___  >__|    \_/  \___  >__|   
     \/|__|           \/     \/                 \/      
		
		Start api server ...`,
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	}
	startCmd = &cobra.Command{
		Use:     "start",
		Short:   "Start api server",
		Example: "web-api start -c config/conf.toml",
		RunE:    start,
	}
)

func init()  {
	cobra.OnInitialize(config.Init)

	startCmd.PersistentFlags().StringVarP(&bootstrap.Args.ConfigFile, "config", "c", "config/conf.toml",
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
