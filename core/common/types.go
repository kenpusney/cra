package common

import "github.com/kenpusney/cra/core/contract"

type Opts struct {
	Endpoint string `arg:"positional"`
	//FollowLocation bool   `arg:"-f" default:"false"`
	Port int `arg:"-p" default:"9511"`
}

type Context interface {
	Register(ty string, strategy Strategy)
	Serve() error
	Proceed(reqItem *contract.RequestItem) *contract.ResponseItem
	Shutdown()
}

type Strategy func(craRequest *contract.Request, context Context, completer contract.ResponseCompleter)

type Config struct {
	BypassedHeaders *[]string `config:"bypassed_headers"`
}
