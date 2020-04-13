package main

import (
	"github.com/kenpusney/cra/core"
)

func main() {

	context := core.NewContext(&core.Opts{
		FollowLocation: false,
	})

	context.Endpoint = "http://localhost:8080"

	context.Serve()
}
