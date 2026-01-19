package cmd

import (
	"fmt"
	"kefir/pkg/image_list"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type GlobalConfig struct {
	Images image_list.Config `json:"images" yaml:"images" mapstructure:"images"`
}

var (
	cfgFile string
	cfg     = GlobalConfig{}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmd.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kefir")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
