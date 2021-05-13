package nsq

import (
	"fmt"
	"reflect"

	proto "github.com/bcowtech/structprototype"
	"github.com/bcowtech/structprototype/reflectutil"
)

type messageHandlerBindingProvider struct {
	data map[string]reflect.Value
}

func (p *messageHandlerBindingProvider) BeforeBind(context *proto.PrototypeContext) error {
	return nil
}

func (p *messageHandlerBindingProvider) BindField(field proto.PrototypeField, rv reflect.Value) error {
	if v, ok := p.data[field.Name()]; ok {
		if !rv.IsValid() {
			return fmt.Errorf("specifiec argument 'rv' is invalid")
		}

		rv = reflectutil.AssignZero(rv)
		rv.Set(v.Convert(rv.Type()))
	}
	return nil
}

func (p *messageHandlerBindingProvider) AfterBind(context *proto.PrototypeContext) error {
	return nil
}
