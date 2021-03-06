package test

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/bcowtech/config"
	"github.com/bcowtech/host"
	nsq "github.com/bcowtech/worker-nsq"
)

func Test(t *testing.T) {
	/* like
	$ go run app.go --handler-concurrency "8"
	*/
	initializeArgs()

	app := MockApp{}
	starter := nsq.Startup(&app,
		[]host.Middleware{
			nsq.UseMessageHandler(&MockMessageHandler{}),
		}...).
		ConfigureConfiguration(func(service *config.ConfigurationService) {
			service.
				LoadEnvironmentVariables("").
				LoadYamlFile("config.yaml").
				LoadCommandArguments()
		})
	// starter.Run()

	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := starter.Start(startCtx); err != nil {
		t.Error(err)
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := starter.Stop(stopCtx); err != nil {
		t.Error(err)
	}

	// assert app.Config
	{
		conf := app.Config
		if conf.NsqConnectionTarget != nsq.NsqLookupd {
			t.Errorf("assert 'Config.NsqConnectionTarget':: expected '%v', got '%v'", nsq.NsqLookupd, conf.NsqConnectionTarget)
		}
		if 0 == len(conf.NsqAddresses) {
			t.Errorf("assert 'Config.NsqAddresses':: should not be empty")
		}
		if conf.Topic != "mytopic" {
			t.Errorf("assert 'Config.Topic':: expected '%v', got '%v'", "mytopic", conf.Topic)
		}
		if conf.Channel != "worker-nsq-demo" {
			t.Errorf("assert 'Config.Channel':: expected '%v', got '%v'", "worker-nsq-demo", conf.Channel)
		}
		if conf.HandlerConcurrency != 8 {
			t.Errorf("assert 'Config.PollingTimeoutMs':: expected '%v', got '%v'", 8, conf.HandlerConcurrency)
		}
	}
}

func initializeArgs() {
	os.Args = []string{"example",
		"--handler-concurrency", "8"}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
