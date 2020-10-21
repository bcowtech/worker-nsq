package nsq

import (
	"fmt"
)

type (
	connectToNSQFunc func(consumer *Consumer, addresses []string) error

	nsqConnectionFuncProvider map[string]connectToNSQFunc
)

var nsqConnectionFactory = nsqConnectionFuncProvider(
	map[string]connectToNSQFunc{
		Nsqd:       connectToNSQDFunc,
		NsqLookupd: connectToNSQLookupdFunc,
	})

func (f nsqConnectionFuncProvider) getInstance(connectionTarget string) (connectToNSQFunc, error) {
	fn, ok := f[connectionTarget]
	if !ok {
		return nil, fmt.Errorf("unknown target '%v'", connectionTarget)
	}
	return fn, nil
}

// Connect to NSQ using nsqd addresses.
func connectToNSQDFunc(consumer *Consumer, addresses []string) error {
	return consumer.ConnectToNSQDs(addresses)
}

// Connect to NSQ using nsqlookupd addresses.
func connectToNSQLookupdFunc(consumer *Consumer, addresses []string) error {
	return consumer.ConnectToNSQLookupds(addresses)
}
