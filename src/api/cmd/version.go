package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/moocss/go-webserver/src/api/pkg/version"
)

// Command returns a command used to print version information.
func versionCmd() *cobra.Command {
	var short bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints out build version information",
		Run: func(cmd *cobra.Command, args []string) {
			if short {
				fmt.Println(version.Info)
			} else {
				fmt.Println(version.Info.LongForm())
			}
		},
	}
	cmd.PersistentFlags().BoolVarP(&short, "version", "v", short, "Displays a short form of the version information")
	return cmd
}