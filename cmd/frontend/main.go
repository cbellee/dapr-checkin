package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cbellee/dapr-checkin/cmd/models"
	"github.com/dapr/go-sdk/service/common"

	dapr "github.com/dapr/go-sdk/client"
	daprd "github.com/dapr/go-sdk/service/http"
)

var (
	version     = "0.0.1"
	serviceName = "frontend"
	servicePort = "8000"
	bindingName = "servicebus.queue.binding"
	queueName   = "checkin"
	logger      = log.New(os.Stdout, "", 0)
)

func main() {
	logger.Printf("Starting service : %v v%v...", serviceName, version)

	port := fmt.Sprintf(":%s", servicePort)
	server := daprd.NewService(port)

	if err := server.AddServiceInvocationHandler("/checkin", checkinHandler); err != nil {
		logger.Panicf("Failed to add service invocation handler '/checkin' : %s", err)
	} else {
		logger.Printf("Invocation handler for service '%v' added successfully!", serviceName)
	}

	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("error listening: %v", err)
	}
}

func checkinHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	if in == nil {
		err = errors.New("invocation parameter required")
		return
	}

	client, err := dapr.NewClient()
	if err != nil {
		logger.Panicf("Failed to create Dapr client: %s", err)
	}

	logger.Printf("echo - ContentType:%s, Verb:%s, QueryString:%s, %s", in.ContentType, in.Verb, in.QueryString, in.Data)

	var c models.Checkin
	err = json.Unmarshal(in.Data, &c)
	if err != nil {
		logger.Print(err.Error())
	}

	br := &dapr.InvokeBindingRequest{
		Name:      bindingName,
		Operation: "create",
		Data:      in.Data,
	}

	if err := client.InvokeOutputBinding(ctx, br); err != nil {
		logger.Panicf("Failed to send event to queue '%s' : %s", queueName, err)
		return nil, err
	} else {
		logger.Printf("Successfully sent event to queue %s", queueName)
	}

	return out, nil
}
