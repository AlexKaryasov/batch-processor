package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type MakeNextBatchTestCase struct {
	NumberOfElements int
	BatchSize        int
	StaringIndex     int
	ExpectedSize     int
}

type BatchProcessTestCase struct {
	MaxBatchSize     int
	FactualBatchSize int
	ExpectedError    error
}

const DefaultTimeout = 1

func TestMakeNextBatch(t *testing.T) {
	testCases := []MakeNextBatchTestCase{
		{11, 10, 0, 10},
		{11, 10, 10, 1},
		{7, 10, 0, 7},
		{0, 10, 0, 0},
	}

	for i, tc := range testCases {
		testCase := tc
		t.Run(fmt.Sprintf("test-case-%d", i), func(t *testing.T) {
			items := make([]Item, testCase.NumberOfElements)
			for i := range items {
				items[i] = struct{}{}
			}
			batchService := BatchService{len(items), testCase.BatchSize, DefaultTimeout}
			out := MakeNextBatch(&items, testCase.StaringIndex, batchService)
			if len(out) != testCase.ExpectedSize {
				t.Errorf("MakeNextBatch returned length %d, expected %d", len(out), testCase.ExpectedSize)
			}
		})
	}
}

func TestProcess(t *testing.T) {
	testCases := []BatchProcessTestCase{
		{10, 11, ErrBlocked},
		{10, 0, ErrNothingToProcess},
		{10, 10, nil},
		{10, 7, nil},
	}

	for i, tc := range testCases {
		testCase := tc
		t.Run(fmt.Sprintf("test-case-%d", i), func(t *testing.T) {
			items := make([]Item, testCase.FactualBatchSize)
			for i := range items {
				items[i] = struct{}{}
			}
			batchService := BatchService{len(items), testCase.MaxBatchSize, DefaultTimeout}
			timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(DefaultTimeout)*time.Second)
			defer cancel()
			out := batchService.Process(timeoutCtx, items)
			if out != testCase.ExpectedError {
				t.Errorf("expected %s error, got %s", testCase.ExpectedError, out)
			}
		})
	}
}
