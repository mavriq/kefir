package cmd

import (
	"kefir/pkg/image_list"
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
