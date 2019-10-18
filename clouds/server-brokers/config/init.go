package config

import "github.com/spf13/viper"

var Cfg *viper.Viper

func init() {
	Cfg = viper.New()

	Cfg.SetConfigFile("")
}
