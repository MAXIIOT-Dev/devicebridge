package cmd

import (
	"fmt"

	"github.com/maxiiot/vbaseBridge/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "print current config of vbaseBridge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Cfg)
	},
}
