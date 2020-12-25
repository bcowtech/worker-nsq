package nsq

import (
	"context"
	"log"
	"sync"

	nsq "github.com/nsqio/go-nsq"
)

type (
	Consumer           = nsq.Consumer
	Config             = nsq.Config
	Message            = nsq.Message
	MessageHandler     = nsq.Handler
	MessageHandlerFunc = nsq.HandlerFunc
)

type Worker struct {
	NsqConnectionTarget string // nsqd, nsqlookupd
	NsqAddresses        []string
	Topic               string
	Channel             string
	HandlerConcurrency  int
	Config              *Config

	messageHandler MessageHandler

	consumer         *Consumer
	consumerSyncOnce sync.Once
	wg               sync.WaitGroup
}

func (w *Worker) Start(ctx context.Context) {
	if w.messageHandler == nil {
		log.Fatalln("[bcow-go/worker-nsq] %% Error: the MessageHandler is nil. Using UseMessagehandler() register one")
	}

	w.consumerSyncOnce.Do(func() {
		w.consumer = w.createConsumer()
	})

	q := w.consumer

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
	connectToNSQ, err := connectionProvider.getInstance(w.NsqConnectionTarget)
	if err != nil {
		log.Fatalf("[bcow-go/worker-nsq] cannot connect to nsq. %v\n", err)
	}
	connectToNSQ(q, w.NsqAddresses)

	// start listening
	log.Printf("[bcow-go/worker-nsq] %s/%s started\n", w.Topic, w.Channel)
	<-q.StopChan
}

func (w *Worker) Stop(ctx context.Context) error {
	log.Printf("[bcow-go/worker-nsq] %% Stop\n")

	if w.consumer != nil {
		w.wg.Wait()
		w.consumer.Stop()
	}
	return nil
}

func (w *Worker) processMessage(m *nsq.Message) error {
	w.wg.Add(1)
	defer func() {
		w.wg.Done()
	}()

	if w.messageHandler != nil {
		return w.messageHandler.HandleMessage(m)
	} else {

	}
	return nil
}

func (w *Worker) createConsumer() *Consumer {
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
	return q
}
