package main

import (
	"log"
	"net/http"
	"os"

	"./spec"
	"./impl"
	"github.com/cbellee/dapr-checkin/cmd/frontend/spec"
	"github.com/dapr/go-sdk/service/common"
)

var (
	version     = "0.0.1"
	buildInfo   = "No build details"
	serviceName = "front-end"
	servicePort = "8001"
	/* cosmosDbName          = "punters"
	cosmosDbContainerName = "default"
	cosmosDbKey           = os.Getenv("COSMOS_DB_KEY")
	cosmosDbURL           = os.Getenv("COSMOS_DB_URL") */
	pubSubName      = "messagebus"
	secretStoreName = "azurekeyvault"
	logger          = log.New(os.Stdout, "", 0)
)

var sub = &common.Subscription{
	PubsubName: pubSubName,
	Topic:      "checkin",
	Route:      "/",
}