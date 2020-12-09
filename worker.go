package nsq

import (
	"context"
	"log"
	"sync"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

const (
	DEFAULT_CONFIG_MAX_IN_FLIGHT         = 8
	DEFAULT_CONFIG_HEARTBEAT_INTERVAL    = 10
	DEFAULT_CONFIG_DEFAULT_REQUEUE_DELAY = 0
	DEFAULT_CONFIG_MAX_BACKOFF_DURATION  = time.Millisecond * 50

	DEFAULT_HANDLER_CONCURRENCY = 12

	Nsqd       = "nsqd"
	NsqLookupd = "nsqlookupd"
)

type (
	Consumer           = nsq.Consumer
	Config             = nsq.Config
	Message            = nsq.Message
	MessageHandlerFunc = nsq.HandlerFunc
)

type Worker struct {
	Consumer            *Consumer
	NsqConnectionTarget string
	NsqAddresses        []string
	Topic               string
	Channel             string
	HandlerConcurrency  int
	Config              *Config
	MessageHandler      MessageHandlerFunc

	consumerSyncOnce sync.Once
	wg               sync.WaitGroup
}

func (w *Worker) Start(ctx context.Context) {
	if w.MessageHandler == nil {
		log.Fatalln("[bcow-go/worker-nsq] %% Error: the MessageHandler cannot be nil.")
	}

	if w.Consumer == nil {
		w.consumerSyncOnce.Do(func() {
			var config *Config = w.Config

			if config == nil {
				c := nsq.NewConfig()
				{
					c.MaxInFlight = DEFAULT_CONFIG_MAX_IN_FLIGHT
					c.HeartbeatInterval = DEFAULT_CONFIG_HEARTBEAT_INTERVAL
					c.DefaultRequeueDelay = DEFAULT_CONFIG_DEFAULT_REQUEUE_DELAY
					c.MaxBackoffDuration = DEFAULT_CONFIG_MAX_BACKOFF_DURATION
				}
				// export
				config = c
			}

			q, err := nsq.NewConsumer(w.Topic, w.Channel, config)
			if err != nil {
				log.Fatalf("[bcow-go/worker-nsq]  %% Error: cannot connect to nsq. %v\n", err)
			}

			w.Consumer = q
		})
	}

	q := w.Consumer

	defer func() {
		q.Stop()
		log.Printf("[bcow-go/worker-nsq] Stats: %+v\n", q.Stats())
		log.Printf("[bcow-go/worker-nsq] IsStarved: %+v\n", q.IsStarved())
	}()

	// bind the MessageHandler
	q.AddConcurrentHandlers(
		nsq.HandlerFunc(w.processMessage),
		w.HandlerConcurrency)

	// get connection and connect to nsqd or lookupd
	connectToNSQ, err := nsqConnectionFactory.getInstance(w.NsqConnectionTarget)
	if err != nil {
		log.Fatalf("[bcow-go/worker-nsq] cannot connect to nsq. %v\n", err)
	}
	connectToNSQ(q, w.NsqAddresses)

	// start listening
	log.Printf("[bcow-go/worker-nsq] %s/%s started\n", w.Topic, w.Channel)
	<-q.StopChan
}

func (w *Worker) Stop(ctx context.Context) error {
	w.wg.Wait()
	w.Consumer.Stop()
	return nil
}

func (w *Worker) processMessage(m *nsq.Message) error {
	w.wg.Add(1)
	defer func() {
		w.wg.Done()
	}()

	if w.MessageHandler != nil {
		return w.MessageHandler(m)
	}
	return nil
}
