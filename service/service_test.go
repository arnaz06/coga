package service_test

import (
	"context"
	"testing"

	"github.com/arnaz06/coga"
	"github.com/arnaz06/coga/service"
	"github.com/stretchr/testify/require"
)

func TestResolveRollback(t *testing.T) {
	tests := []struct {
		testName                string
		serviceName             string
		transactionList         []coga.TransactionList
		expectedTransactionList []coga.TransactionList
	}{
		{
			testName:    "success rollback shipment service",
			serviceName: "shipment-service",
			transactionList: []coga.TransactionList{
				{
					ServiceName: "coga-service",
					Topic:       "staging-saga-service",
				},
				{
					ServiceName: "order-service",
					Topic:       "staging-saga-service-order",
				},
				{
					ServiceName: "inventory-service",
					Topic:       "staging-saga-service-inventory",
				},
				{
					ServiceName: "shipment-service",
					Topic:       "staging-saga-service-shipment",
				},
			},
			expectedTransactionList: []coga.TransactionList{
				{
					ServiceName: "order-service",
					Topic:       "staging-saga-service-order",
				},
				{
					ServiceName: "inventory-service",
					Topic:       "staging-saga-service-inventory",
				},
			},
		},
		{
			testName:    "success rollback inventory service",
			serviceName: "inventory-service",
			transactionList: []coga.TransactionList{
				{
					ServiceName: "coga-service",
					Topic:       "staging-saga-service",
				},
				{
					ServiceName: "order-service",
					Topic:       "staging-saga-service-order",
				},
				{
					ServiceName: "inventory-service",
					Topic:       "staging-saga-service-inventory",
				},
				{
					ServiceName: "shipment-service",
					Topic:       "staging-saga-service-shipment",
				},
			},
			expectedTransactionList: []coga.TransactionList{
				{
					ServiceName: "order-service",
					Topic:       "staging-saga-service-order",
				},
			},
		},
		{
			testName:    "success rollback order service",
			serviceName: "order-service",
			transactionList: []coga.TransactionList{
				{
					ServiceName: "coga-service",
					Topic:       "staging-saga-service",
				},
				{
					ServiceName: "order-service",
					Topic:       "staging-saga-service-order",
				},
				{
					ServiceName: "inventory-service",
					Topic:       "staging-saga-service-inventory",
				},
				{
					ServiceName: "shipment-service",
					Topic:       "staging-saga-service-shipment",
				},
			},
			expectedTransactionList: []coga.TransactionList{},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			mh := service.NewMessaheHandler(context.Background(), test.transactionList)
			res := mh.ResolveRollback(test.transactionList, test.serviceName)
			require.Equal(t, test.expectedTransactionList, res)
		})
	}
}

func TestResolveNextService(t *testing.T) {
	tests := []struct {
		testName            string
		serviceName         string
		transactionList     []coga.TransactionList
		expectedNextService string
		expectedNextTopic   string
	}{
		{
			testName:    "success",
			serviceName: "order-service",
			transactionList: []coga.TransactionList{
				{
					ServiceName: "order-service",
					Topic:       "staging-saga-service-order",
				},
				{
					ServiceName: "inventory-service",
					Topic:       "staging-saga-service-inventory",
				},
				{
					ServiceName: "shipment-service",
					Topic:       "staging-saga-service-shipment",
				},
			},
			expectedNextService: "inventory-service",
			expectedNextTopic:   "staging-saga-service-inventory",
		},
		{
			testName:    "success inventory service",
			serviceName: "inventory-service",
			transactionList: []coga.TransactionList{
				{
					ServiceName: "order-service",
					Topic:       "staging-saga-service-order",
				},
				{
					ServiceName: "inventory-service",
					Topic:       "staging-saga-service-inventory",
				},
				{
					ServiceName: "shipment-service",
					Topic:       "staging-saga-service-shipment",
				},
			},
			expectedNextService: "shipment-service",
			expectedNextTopic:   "staging-saga-service-shipment",
		},
		{
			testName:    "last transaction list",
			serviceName: "shipment-service",
			transactionList: []coga.TransactionList{
				{
					ServiceName: "order-service",
					Topic:       "staging-saga-service-order",
				},
				{
					ServiceName: "inventory-service",
					Topic:       "staging-saga-service-inventory",
				},
				{
					ServiceName: "shipment-service",
					Topic:       "staging-saga-service-shipment",
				},
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
