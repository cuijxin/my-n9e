package config

import (
	"bytes"
	"fmt"

	"github.com/cuijxin/my-n9e/models"
	"github.com/cuijxin/my-n9e/pkg/iconf"
	"github.com/cuijxin/my-n9e/pkg/ilog"
	"github.com/ory/viper"
	"github.com/toolkits/pkg/file"
)

type ConfigStruct struct {
	Logger ilog.Config         `yaml:"logger"`
	HTTP   httpSection         `yaml:"http"`
	MySQL  models.MysqlSection `yaml:"mysql"`
}

type httpSection struct {
	Mode           string `yaml:"mode"`
	Access         bool   `yaml:"access"`
	Listen         string `yaml:"listen"`
	Pprof          bool   `yaml:"pprof"`
	CookieName     string `yaml:"cookieName"`
	CookieDomain   string `yaml:"cookieDomain"`
	CookieSecure   bool   `yaml:"cookieSecure"`
	CookieHttpOnly bool   `yaml:"cookieHttpOnly"`
	CookieMaxAge   int    `yaml:"cookieMaxAge"`
	CookieSecret   string `yaml:"cookieSecret"`
	CsrfSecret     string `yaml:"csrfSecret"`
}

var Config *ConfigStruct

func Parse() error {
	ymlFile := iconf.GetYmlFile("server")
	if ymlFile == "" {
		return fmt.Errorf("configuration file of server not found")
	}

	bs, err := file.ReadBytes(ymlFile)
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", ymlFile, err)
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(bs))
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", ymlFile, err)
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", ymlFile, err)
	}

	fmt.Println("config.file:", ymlFile)

	return nil
}
