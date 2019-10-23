package config

import (
	"flag"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/spf13/viper"
)

var Cfg *viper.Viper
var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "f", "G:\\go_workspace\\GOPATH\\src\\tcpx\\projects\\jelly\\config\\pro.yaml", "-f <cfg/file/path>")
	flag.Parse()

	Cfg = viper.New()
	Cfg.SetConfigFile(configFilePath)
	if e := Cfg.ReadInConfig(); e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		return
	}
}
