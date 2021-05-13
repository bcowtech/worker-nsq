package nsq

import (
	"github.com/bcowtech/host"
)

func Startup(app interface{}, middlewares ...host.Middleware) *host.Starter {
	starter := host.Startup(app, middlewares...)
	// set options

	// init HostProvider
	host.SetupHostProvider(starter, nsqHostProvider)

	return starter
}
