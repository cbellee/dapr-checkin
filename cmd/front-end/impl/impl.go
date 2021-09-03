package impl

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cbellee/dapr-checkin/cmd/front-end/spec"
	"github.com/cbellee/dapr-checkin/pkg/helper"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
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
	defer helper.TimeTrack(time.Now(), "NewService()")

	client, err := dapr.NewClient()
	if err != nil {
		logger.Panicf("Failed to create Dapr client: %s", err)
		return nil
	}

	port := fmt.Sprintf(":%s", servicePort)
	server := daprd.NewService(port)
	if err != nil {
		logger.Panicf("Failed to create Dapr server: %s", err)
	}

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
		logger.Panicf("Failed to start service : %s", err)
		return err
	}

	return nil
}

// PublishEvent publishes an event to a message bus topic
func (s *Service) PublishEvent(pubsubName string, topicName string, data []byte) error {
	defer helper.TimeTrack(time.Now(), "PublshEvent()")

	ctx := context.Background()
	if err := s.client.PublishEvent(ctx, pubsubName, topicName, data); err != nil {
		logger.Panicf("Failed to publish event to topic '%s' : %s", topicName, err)
		return err
	} else {
		logger.Printf("Published event to topic '%s' on pubsub '%s'", topicName, pubsubName)
	}

	return nil
}
