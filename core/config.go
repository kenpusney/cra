package core

import (
	"github.com/yangeagle/config"
)

type Config struct {
	BypassedHeaders *[]string `config:"bypassed_headers"`
}

func LoadConfig() *Config {
	confParser := config.NewConfig()

	err := confParser.ParseFile("cra.conf")

	if err != nil {
		return nil
	}

	conf := new(Config)

	err = confParser.Unmarshal(&conf)
	if err != nil {
		return nil
	}

	return conf
}
