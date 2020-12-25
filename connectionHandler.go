package nsq

import (
	"fmt"
)

type (
	connectToNSQFunc func(consumer *Consumer, addresses []string) error

	connectionFuncProvider map[string]connectToNSQFunc
)

var connectionProvider = connectionFuncProvider(
	map[string]connectToNSQFunc{
		Nsqd:       connectToNSQD,
		NsqLookupd: connectToNSQLookupd,
	})

func (f connectionFuncProvider) getInstance(connectionTarget string) (connectToNSQFunc, error) {
	fn, ok := f[connectionTarget]
	if !ok {
		return nil, fmt.Errorf("unknown target '%v'", connectionTarget)
	}
	return fn, nil
}

// Connect to NSQ using nsqd addresses.
func connectToNSQD(consumer *Consumer, addresses []string) error {
	return consumer.ConnectToNSQDs(addresses)
}

// Connect to NSQ using nsqlookupd addresses.
func connectToNSQLookupd(consumer *Consumer, addresses []string) error {
	return consumer.ConnectToNSQLookupds(addresses)
}
