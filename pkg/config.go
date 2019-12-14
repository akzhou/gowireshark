/**
 * @Author: Administrator
 * @Description:
 * @File:  config
 * @Version: 1.0.0
 * @Date: 2019/12/12 13:52
 */

package pkg

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"runtime"
)

type WireSharkConfig struct {
	UrlFlag        string
	FileServerPort uint16
}

var (
	wireSharkCfg WireSharkConfig
	confPath     string
)

func init() {
	var tomlPath string
	if runtime.GOOS == `windows` {
		tomlPath = "e:/xinxinserver/config/gowireshark.toml"
	} else {
		tomlPath = "/config/gowireshark.toml"
	}
	flag.StringVar(&confPath, "conf", tomlPath, "config path")

	viper.SetConfigFile(confPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&wireSharkCfg); err != nil {
		panic(err)
	}

	if wireSharkCfg.UrlFlag == "" {
		panic(fmt.Errorf("未配置UrlFlag"))
	}

	if wireSharkCfg.FileServerPort == 0 {
		panic(fmt.Errorf("未配置FileServerPort"))
	}
}
