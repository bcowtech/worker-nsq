package nsq

import (
	"gitlab.bcowtech.de/bcow-go/host"
)

var (
	defaultHostProvider = &hostProvider{}
)

func Startup(app interface{}, middlewares ...host.Middleware) *host.Starter {
	starter := host.Startup(app, middlewares...)
	// set options

	// init HostProvider
	starter.HostProviderWrapper.Wrap(defaultHostProvider)

	return starter
}
