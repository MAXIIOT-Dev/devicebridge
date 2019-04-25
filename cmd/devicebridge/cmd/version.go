package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version of vbaseBridge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
