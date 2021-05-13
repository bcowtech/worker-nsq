package nsq

import (
	"time"

	"github.com/bcowtech/host"
	"github.com/nsqio/go-nsq"
)

const (
	DEFAULT_CONFIG_MAX_IN_FLIGHT         = 8
	DEFAULT_CONFIG_HEARTBEAT_INTERVAL    = 10
	DEFAULT_CONFIG_DEFAULT_REQUEUE_DELAY = 0
	DEFAULT_CONFIG_MAX_BACKOFF_DURATION  = time.Millisecond * 50

	DEFAULT_HANDLER_CONCURRENCY = 12

	Nsqd       = "nsqd"
	NsqLookupd = "nsqlookupd"

	appHostFieldName            = host.AppHostFieldName
	appConfigFieldName          = host.AppConfigFieldName
	appServiceProviderFieldName = host.AppServiceProviderFieldName
	componentInitMethodName     = host.ComponentInitMethodName
)

var (
	nsqHostProvider = &hostProvider{}
)

func NewMessage(id MessageID, body []byte) *Message {
	return nsq.NewMessage(id, body)
}
