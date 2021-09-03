package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cbellee/dapr-checkin/cmd/front-end/impl"
	"github.com/cbellee/dapr-checkin/cmd/front-end/spec"
	"github.com/dapr/go-sdk/service/common"
)

var (
	version         = "0.0.1"
	buildInfo       = "No build details"
	serviceName     = "back-end"
	servicePort     = "8001"
	messageBusName  = ""
	storeName       = "checkin-statestore"
	pubSubName      = "messages"
	topicName = "checkinEvents"
	secretStoreName = "azurekeyvault"
	logger          = log.New(os.Stdout, "", 0)
)

var components = spec.DaprComponents{
	MessageBusName: messageBusName,
	TopicName:      topicName,
}

var sub = &common.Subscription{
	PubsubName: pubSubName,
	Topic:      "checkin",
	Route:      "/checkins",
	Metadata:   nil,
}

// API type
type API struct {
	service spec.Service
}

func main() {
	logger.Printf("### Dapr: %v v%v starting...", serviceName, version)

	api := API{
		impl.NewService(serviceName, servicePort, components),
	}

	if err := api.service.AddTopicHandler("/add", api.service.AddBet); err != nil {
		logger.Fatalf("error adding 'AddBet' invocation handler: %v", err)
	}

	if err := api.service.StartService(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("error: %v", err)
	}
}
