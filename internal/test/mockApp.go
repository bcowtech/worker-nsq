package test

import (
	nsq "github.com/bcowtech/worker-nsq"
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
		ResouceName string
	}
)

func (provider *ServiceProvider) Init(conf *Config) {
	provider.ResouceName = "demo resource"
}

func (h *Host) Init(conf *Config) {
	h.NsqConnectionTarget = conf.NsqConnectionTarget
	h.NsqAddresses = conf.NsqAddresses
	h.Topic = conf.Topic
	h.Channel = conf.Channel
	h.HandlerConcurrency = conf.HandlerConcurrency
}
