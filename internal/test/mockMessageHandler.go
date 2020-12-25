package test

import (
	"fmt"
	"time"

	nsq "gitlab.bcowtech.de/bcow-go/worker-nsq"
)

type MockMessageHandler struct {
	ServiceProvider *ServiceProvider
}

func (h *MockMessageHandler) Init() {
	fmt.Println("MockMessageHandler.Init()")
	fmt.Printf("MockMessageHandler.ServiceProvider: %+v\n", h.ServiceProvider)
}

func (h *MockMessageHandler) HandleMessage(message *nsq.Message) error {
	if len(message.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return nil
	}

	fmt.Println("======= message start")
	fmt.Println(string(message.Body))
	time.Sleep(time.Second * 4)
	fmt.Println("======= message end")

	return nil
}
