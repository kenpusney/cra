package boot

import (
	"github.com/kenpusney/cra/core/common"
	"github.com/kenpusney/cra/core/context"
	"github.com/kenpusney/cra/core/strategy"
	"github.com/yangeagle/config"
)

func NewContext(opts *common.Opts) common.Context {

	ctx := context.MakeContext(opts, LoadConfig())

	ctx.Register("seq", strategy.Sequential)
	ctx.Register("con", strategy.Concurrent)
	ctx.Register("cascaded", strategy.Cascaded)
	ctx.Register("batch", strategy.Batch)

	return ctx
}

func LoadConfig() *common.Config {
	confParser := config.NewConfig()

	err := confParser.ParseFile("cra.conf")

	if err != nil {
		return nil
	}

	conf := new(common.Config)

	err = confParser.Unmarshal(&conf)
	if err != nil {
		return nil
	}

	return conf
}
