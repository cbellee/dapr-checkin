package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"../spec"

	"../../../pkg/helper"
	dapr "github.com/dapr/go-sdk/client"
	common "github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

var (
	logger = log.New(os.Stdout, "", 0)
)

// Service implements a dapr service and client
type Service struct {
	client     dapr.Client
	server     common.Service
	components spec.DaprComponents
}

// NewService creates a new instance of the Service
func NewService(serviceName string, servicePort string, components spec.DaprComponents) *Service {
	client, err := dapr.NewClient()
	if err != nil {
		logger.Panicf("Failed to create Dapr client: %s", err)
		return nil
	}

	port := fmt.Sprintf(":%s", servicePort)
	server := daprd.NewService(port)

	service := &Service{
		client,
		server,
		components,
	}

	return service
}

// StartService starts the http server
func (s *Service) StartService() error {
	defer helper.TimeTrack(time.Now(), "StartService()")

	err := s.server.Start()
	if err != nil {
		return err
	}
	return nil
}

func eventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	log.Printf("event - PubsubName: %s, Topic: %s, ID: %s, Data: %s", e.PubsubName, e.Topic, e.ID, e.Data)
	return false, nil
}

// AddTopicHandler wires up a new topic event handler
func (s *Service) AddTopicHandler(sub *common.Subscription, handler func(context.Context, *common.TopicEvent) error) error {
	defer helper.TimeTrack(time.Now(), "AddTopicHandler()")

	if err := s.server.AddTopicEventHandler(sub, eventHandler); err != nil {
		logger.Fatalf("Error adding topic event handler: %s", err)
		return err
	}

	return nil
}

// TopicHandler listens for new checkins from a Pub/Sub Topic
func (s *Service) TopicHandler(ctx context.Context, e *common.TopicEvent) error {
	defer helper.TimeTrack(time.Now(), "TopicHandler()")

	d, err := json.Marshal(e.Data)
	if err != nil {
		logger.Fatal(err)
		return err
	}

	var result spec.Result
	if err := json.Unmarshal(d, &result); err != nil {
		logger.Fatal(err)
	}

	return nil
}

func (s *Service) SaveState(ctx context.Context, storeName string, key string, data []byte) error {
	defer helper.TimeTrack(time.Now(), "SaveState()")

	if err := s.client.SaveState(ctx, storeName, key, data); err != nil {
		logger.Panicf("Failed to save state to state store '%s' : %s", storeName, err)
		return err
	} else {
		logger.Printf("Saved state to statestore '%s' with key '%s'", storeName, key)
	}

	return nil
}
