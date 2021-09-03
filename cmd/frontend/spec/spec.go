package spec

import (
	"context"
	"time"

	"github.com/dapr/go-sdk/service/common"
)

// Checkin is a struct
type Checkin struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	LocationID string    `json:"location_id"`
	TimeStamp  time.Time `json:"timestamp"`
}

// DaprComponents is a struct
type DaprComponents struct {
	TopicName      string `json:"topic_name"`
	MessageBusName string `json:"messagebus_name"`
	BindingName    string `json:"binding_name"`
}

// Service defines the behaviours needed to interact with the service
type Service interface {
	AddTopicHandler(sub *common.Subscription, fn func(context.Context, *common.TopicEvent) error) error
	AddServiceHandler(name string, fn func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error)) error
	PublishEvent(pubsubName string, topicName string, data []byte) error
	StartService() error
}
