package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Command returns a command used to print version information.
func Command() *cobra.Command {
	var short bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints out build version information",
		Run: func(cmd *cobra.Command, args []string) {
			if short {
				fmt.Println(Info)
			} else {
				fmt.Println(Info.LongForm())
			}
		},
	}
	cmd.PersistentFlags().BoolVarP(&short, "version", "v", short, "Displays a short form of the version information")
	return cmd
}
