package cmd

import (
	"os"

	"github.com/maxiiot/vbaseBridge/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var version string

var rootCmd = &cobra.Command{
	Use:   "vbase-bridge",
	Short: "run vbase-bridge server",
	RunE:  run,
}

// Execute start server
func Execute(v string) {
	version = v
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to configuration file (optional)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
}

func initConfig() {
	if cfgFile != "" {
		file, err := os.Open(cfgFile)
		if err != nil {
			log.WithError(err).WithField("config", cfgFile).Fatal("error loading config file")
		}
		viper.SetConfigType("toml")
		if err := viper.ReadConfig(file); err != nil {
			log.WithError(err).WithField("config", cfgFile).Fatal("error loading config file")
		}
	} else {
		viper.SetConfigName("app")
		viper.AddConfigPath("./")
		viper.AddConfigPath("$HOME/.config/vbaseBridge")
		viper.AddConfigPath("/etc/vbaseBridge/")

		if err := viper.ReadInConfig(); err != nil {
			log.WithError(err).Fatal("read config file error.")
		}
	}

	if err := viper.Unmarshal(&config.Cfg); err != nil {
		log.WithError(err).Fatal("error unmarshal config file.")
	}
}
