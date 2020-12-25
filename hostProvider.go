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

func (p *hostProvider) Init(h host.Host, ctx *host.Context) {
}

func (p *hostProvider) PostLoadMiddleware(h host.Host, ctx *host.Context) {
}

func (p *hostProvider) Emit(rv reflect.Value) host.Host {
	rvHost := reflect.NewAt(typeOfHost, unsafe.Pointer(rv.Pointer()))
	v, ok := rvHost.Interface().(host.Host)
	if ok {
		return v
	}
	return nil
}

func (p *hostProvider) asNsqWorker(rv reflect.Value) *Worker {
	return reflect.NewAt(typeOfHost, unsafe.Pointer(rv.Pointer())).
		Interface().(*Worker)
}
