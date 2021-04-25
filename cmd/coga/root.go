package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"cloud.google.com/go/pubsub"
	"github.com/alexsasharegan/dotenv"
	"github.com/arnaz06/coga"
	_pubsub "github.com/arnaz06/coga/internal/pubsub"
	service "github.com/arnaz06/coga/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var rootCmd = &cobra.Command{
	Use:   "coga",
	Short: "co-saga",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	cobra.OnInitialize(initApp)
}

// Execute the main function
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func initApp() {
	//--Initiating coga service--//
	log.Info("Starting SAGA...")
	err := dotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err)
	}
	projectID := os.Getenv("PUBSUB_PROJECT_ID")
	credentialFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	maxOutstandingBytes, err := strconv.ParseInt(os.Getenv("PUBSUB_MAX_OUTSTANDING_BYTES"), 10, 64)
	if err != nil {
		log.Fatal("PUBSUB_MAX_OUTSTANDING_BYTES is not well-set")
	}

	outsandingMessage, err := strconv.ParseInt(os.Getenv("PUBSUB_MAX_OUTSTANDING_MESSAGE"), 10, 64)
	if err != nil {
		log.Fatal("PUBSUB_MAX_OUTSTANDING_MESSAGE is not well-set")
	}

	numWorker, err := strconv.ParseInt(os.Getenv("PUBSUB_NUM_WORKER"), 10, 64)
	if err != nil {
		log.Fatal("PUBSUB_NUM_WORKER is not well-set")
	}

	pubsubClient, err := pubsub.NewClient(context.Background(), projectID, option.WithCredentialsFile(credentialFile))
	if err != nil {
		log.Fatalf("error initializing google cloud pubsub, err:%+v", err)
	}

	subs := pubsubClient.Subscription("staging-subs-saga-service")
	subs.ReceiveSettings.MaxOutstandingMessages = int(outsandingMessage)
	subs.ReceiveSettings.NumGoroutines = int(numWorker)
	subs.ReceiveSettings.MaxOutstandingBytes = int(maxOutstandingBytes)

	transactionFile := os.Getenv("TRANSACTION_CONF")
	jsonFile, err := os.Open(transactionFile)
	if err != nil {
		log.Fatalf("error open transaction list file, err:%+v", err)
	}
	defer jsonFile.Close()

	tlJSON, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("error read transaction list file, err:%+v", err)
	}

	transactionList := []coga.TransactionList{}
	err = json.Unmarshal(tlJSON, &transactionList)
	if err != nil {
		log.Fatalf("error unmarshal transaction list file, err:%+v", err)
	}

	if len(transactionList) == 0 {
		log.Fatal("transaction list can't be empty")
	}

	ctx := context.Background()
	for _, tl := range transactionList {
		fmt.Println("==============")
		fmt.Println("==== Service Name", tl.ServiceName)
		fmt.Println("==== Topic Name", tl.Topic)
		fmt.Println("==============")
		topic := pubsubClient.Topic(tl.Topic)
		publisher := _pubsub.NewEventPublishingService(ctx, topic)
		ctx = context.WithValue(ctx, tl.Topic, publisher)
	}

	cogaHandler := service.NewMessaheHandler(ctx, transactionList)
	messageHandler := _pubsub.NewSubMessageHandler(subs, cogaHandler)

	messageHandler.Pull()
}
