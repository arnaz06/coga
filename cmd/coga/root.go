package main

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
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
	// projectID := os.Getenv("PUBSUB_PROJECT_ID")
	credentialFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	// subscriptionName := os.Getenv("PUBSUB_SUBSCRIPTION_NAME")

	// synchronous, err := strconv.ParseBool(os.Getenv("PUBSUB_SYNCHRONOUS"))
	// if err != nil {
	// 	log.Fatal("PUBSUB_SYNCHRONOUS is not well-set")
	// }

	// maxOutstandingBytes, err := strconv.ParseInt(os.Getenv("PUBSUB_MAX_OUTSTANDING_BYTES"), 10, 64)
	// if err != nil {
	// 	log.Fatal("PUBSUB_MAX_OUTSTANDING_BYTES is not well-set")
	// }

	// outsandingMessage, err := strconv.ParseInt(os.Getenv("PUBSUB_MAX_OUTSTANDING_MESSAGE"), 10, 64)
	// if err != nil {
	// 	log.Fatal("PUBSUB_MAX_OUTSTANDING_MESSAGE is not well-set")
	// }

	// numWorker, err := strconv.ParseInt(os.Getenv("PUBSUB_NUM_WORKER"), 10, 64)
	// if err != nil {
	// 	log.Fatal("PUBSUB_NUM_WORKER is not well-set")
	// }
	fmt.Println(credentialFile + "ulalala")
	pubsubClient, err := pubsub.NewClient(context.Background(), "lemonilo-1487765880311", option.WithCredentialsFile("../../google_service_account.json"))
	if err != nil {
		log.Fatalf("error initializing google cloud pubsub, err:%+v", err)
	}

	subs := pubsubClient.Subscription("staging-subs-saga-service")
	subs.ReceiveSettings.MaxOutstandingMessages = int(100)
	subs.ReceiveSettings.NumGoroutines = int(1)
	// subs.ReceiveSettings.Synchronous = synchronous
	subs.ReceiveSettings.MaxOutstandingBytes = int(100000)

	transactionList := map[string]string{
		"order-service":     "staging-saga-service-order",
		"inventory-service": "staging-saga-service-inventory",
		"shipment-service":  "staging-saga-service-shipment",
	}

	ctx := context.Background()
	for i, tl := range transactionList {
		fmt.Println("==============")
		fmt.Println("====Service Name", i)
		fmt.Println("====Topic Name", tl)
		fmt.Println("==============")
		topic := pubsubClient.Topic(tl)
		publisher := _pubsub.NewEventPublishingService(ctx, topic)
		ctx = context.WithValue(ctx, tl, publisher)
	}

	cogaHandler := service.NewMessaheHandler(ctx, transactionList)
	// subName := pubsubClient.Subscription("staging-subs-saga-service")
	messageHandler := _pubsub.NewSubMessageHandler(subs, cogaHandler)

	messageHandler.Pull()
}
