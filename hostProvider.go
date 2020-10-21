package nsq

import (
	"reflect"
	"unsafe"

	"gitlab.bcowtech.de/bcow-go/host"
)

var (
	typeOfHost = reflect.TypeOf(Worker{})
)

type hostProvider struct{}

func (p *hostProvider) Init(h host.Host, app *host.AppContext) {
}

func (p *hostProvider) ConfigureHostComponent(h host.Host, app *host.AppContext) {
}

func (p *hostProvider) ConvertFromValue(rv reflect.Value) host.Host {
	rvHost := reflect.NewAt(typeOfHost, unsafe.Pointer(rv.Pointer()))
	return rvHost.Interface().(host.Host)
}
