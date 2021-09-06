package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cbellee/dapr-checkin/cmd/models"
	"github.com/dapr/go-sdk/service/common"
	"github.com/google/uuid"

	dapr "github.com/dapr/go-sdk/client"
	daprd "github.com/dapr/go-sdk/service/http"
)

var (
	version          = "0.0.1"
	serviceName      = "backend"
	servicePort      = "8001"
	queueBindingName = "servicebus.queue.binding"
	logger           = log.New(os.Stdout, "", 0)
	storeBindingName = "cosmosdb.store.binding"
)

/*
var sub = &common.Subscription{
	PubsubName: pubSubName,
	Topic:      "checkin",
	Route:      "/events",
}
*/

func main() {
	logger.Printf("### Dapr: %v v%v starting...", serviceName, version)

	port := fmt.Sprintf(":%s", servicePort)
	server := daprd.NewService(port)

	if err := server.AddBindingInvocationHandler(queueBindingName, checkinHandler); err != nil {
		logger.Panicf("Failed to add queue binding invocation handler : %s", err)
	}

	if err := server.Start(); err != nil {
		logger.Panicf("Failed to start service : %s", err)
	}
}

func checkinHandler(ctx context.Context, e *common.BindingEvent) (out []byte, err error) {
	logger.Printf("event - Data: %s, MetaData: %s", e.Data, e.Metadata)

	var checkin models.Checkin
	if err := json.Unmarshal(e.Data, &checkin); err != nil {
		logger.Fatal(err)
		return nil, err
	}

	id := uuid.New()
	checkin.ID = id.String()

	// save to state store
	saveCheckin(ctx, &checkin)

	return out, nil
}

func saveCheckin(ctx context.Context, in *models.Checkin) (retry bool, err error) {

	// create dapr client
	client, err := dapr.NewClient()
	if err != nil {
		logger.Panicf("Failed to create Dapr client: %s", err)
	}

	bytArr, err := json.Marshal(in)
	if err != nil {
		logger.Print(err.Error())
	}

	br := &dapr.InvokeBindingRequest{
		Name:      storeBindingName,
		Data:      bytArr,
		Operation: "create",
	}

	// save message to state store using output binding
	logger.Printf("invoking binding '%s'", storeBindingName)
	err = client.InvokeOutputBinding(ctx, br)
	if err != nil {
		logger.Print(err.Error())
	} else {
		logger.Printf("new checkin with UserID: '%s' LocationID: '%s' CheckinTime: '%d' saved successfully", in.UserID, in.LocationID, in.CheckInTimeStamp)
	}

	return false, nil
}
