package nsq

import (
	"log"
	"reflect"

	"github.com/bcowtech/host"
	proto "github.com/bcowtech/structprototype"
)

func UseMessageHandler(handler MessageHandler) host.Middleware {
	if handler == nil {
		panic("argument 'handler' cannot be nil")
	}

	return &host.GenericMiddleware{
		InitFunc: func(appCtx *host.Context) {
			rvHost := appCtx.HostField()
			worker := nsqHostProvider.asNsqWorker(rvHost)

			err := bindMessageHandler(handler, appCtx)
			if err != nil {
				panic(err)
			}

			worker.messageHandler = handler
		},
	}
}

func bindMessageHandler(v interface{}, appCtx *host.Context) error {
	// populate the ServiceProvider & Config
	provider := &messageHandlerBindingProvider{
		data: map[string]reflect.Value{
			appConfigFieldName:          appCtx.FieldByName(appConfigFieldName),
			appServiceProviderFieldName: appCtx.FieldByName(appServiceProviderFieldName),
		},
	}

	prototype, err := proto.Prototypify(v,
		&proto.PrototypifyConfig{
			BuildValueBinderFunc: proto.BuildNilBinder,
			StructTagResolver:    resolveMessageHandlerTag,
		})
	if err != nil {
		return err
	}

	binder, err := proto.NewPrototypeBinder(prototype, provider)
	if err != nil {
		return err
	}

	err = binder.Bind()
	if err != nil {
		return err
	}

	ctx := proto.PrototypeContext(*prototype)
	rv := ctx.Target()
	if rv.CanAddr() {
		rv = rv.Addr()
		// call MessageHandler.Init()
		fn := rv.MethodByName(componentInitMethodName)
		if fn.IsValid() {
			if fn.Kind() != reflect.Func {
				log.Fatalf("[bcow-go/worker-nsq] cannot find %s.%s() within type %[1]s\n", rv.Type().String(), componentInitMethodName)
			}
			if fn.Type().NumIn() != 0 || fn.Type().NumOut() != 0 {
				log.Fatalf("[bcow-go/worker-nsq] %s.%s() type should be func()\n", rv.Type().String(), componentInitMethodName)
			}
			fn.Call([]reflect.Value(nil))
		}
	}
	return nil
}

func resolveMessageHandlerTag(fieldname, token string) (*proto.StructTag, error) {
	tag := &proto.StructTag{
		Name: fieldname,
	}
	return tag, nil
}
