package service_test

import (
	"context"
	"testing"

	"github.com/arnaz06/coga/service"
	"github.com/stretchr/testify/require"
)

func TestResolveNextService(t *testing.T) {
	tests := []struct {
		testName            string
		serviceName         string
		transactionList     map[string]string
		expectedNextService string
		expectedNextTopic   string
	}{
		{
			testName:    "success",
			serviceName: "order-service",
			transactionList: map[string]string{
				"order-service":     "staging-saga-service-order",
				"inventory-service": "staging-saga-service-inventory",
				"shipment-service":  "staging-saga-service-shipment",
			},
			expectedNextService: "inventory-service",
			expectedNextTopic:   "staging-saga-service-inventory",
		},
		{
			testName:    "success",
			serviceName: "inventory-service",
			transactionList: map[string]string{
				"order-service":     "staging-saga-service-order",
				"inventory-service": "staging-saga-service-inventory",
				"shipment-service":  "staging-saga-service-shipment",
			},
			expectedNextService: "shipment-service",
			expectedNextTopic:   "staging-saga-service-shipment",
		},
		{
			testName:    "last transaction",
			serviceName: "shipment-service",
			transactionList: map[string]string{
				"order-service":     "staging-saga-service-order",
				"inventory-service": "staging-saga-service-inventory",
				"shipment-service":  "staging-saga-service-shipment",
			},
			expectedNextService: "",
			expectedNextTopic:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			mh := service.NewMessaheHandler(context.Background(), test.transactionList)
			nextService, nextTopic := mh.ResolveNextService(test.transactionList, test.serviceName)
			require.Equal(t, test.expectedNextService, nextService)
			require.Equal(t, test.expectedNextTopic, nextTopic)
		})
	}
}
