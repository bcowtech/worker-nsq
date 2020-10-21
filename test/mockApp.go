package test

import (
	"fmt"
	"time"

	nsq "gitlab.bcowtech.de/bcow-go/worker-nsq"
)

type (
	MockApp struct {
		Host            *Host
		Config          *Config
		ServiceProvider *ServiceProvider
	}

	Host nsq.Worker

	Config struct {
		// nsq
		NsqConnectionTarget string   `env:"*NSQ_CONNECTION_TARGET"`
		NsqAddresses        []string `env:"*NSQ_ADDRESSES"`
		Topic               string   `env:"-"                       yaml:"topic"`
		Channel             string   `env:"-"                       yaml:"channel"`
		HandlerConcurrency  int      `env:"-"                       yaml:"handlerConcurrency"    arg:"handler-concurrency"`
	}

	ServiceProvider struct {
	}
)

func (provider *ServiceProvider) Init(conf *Config) {
}

func (h *Host) Init(conf *Config) {
	h.NsqConnectionTarget = conf.NsqConnectionTarget
	h.NsqAddresses = conf.NsqAddresses
	h.Topic = conf.Topic
	h.Channel = conf.Channel
	h.HandlerConcurrency = conf.HandlerConcurrency
	h.MessageHandler = func(m *nsq.Message) error {
		if len(m.Body) == 0 {
			// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
			// In this case, a message with an empty body is simply ignored/discarded.
			return nil
		}

		fmt.Println("======= message start")
		fmt.Println(string(m.Body))
		time.Sleep(time.Duration(4) * time.Second)
		fmt.Println("======= message end")

		return nil
	}
}
