package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrBlocked reports if service is blocked.
var ErrBlocked = errors.New("blocked")
var ErrNothingToProcess = errors.New("nothing to process")

const (
	InputNumberOfElements = "Please input how many elements are you going to process"
	InputBatchSize        = "Please input batch size"
	InputTimeout          = "Please input batch timeout"
)

// Service defines external service that can process batches of items.
type Service interface {
	Process(ctx context.Context, batch Batch) error
}

// BatchService is a struct that encapsulates input parameters
type BatchService struct {
	NumberOfElements int
	BatchSize        int
	Timeout          int
}

func setFieldToBatchService(message string, field *int) {
	fmt.Println(message)
	for {
		_, err := fmt.Scan(field)
		if err != nil {
			fmt.Println(message)
		} else {
			break
		}
	}
}

func MakeNextBatch(items *[]Item, startIndex int, batchService BatchService) []Item {
	if startIndex+batchService.BatchSize >= batchService.NumberOfElements {
		return (*items)[startIndex:]
	}
	return (*items)[startIndex : startIndex+batchService.BatchSize]
}

func (bs BatchService) Process(ctx context.Context, batch Batch) error {
	if len(batch) == 0 {
		return ErrNothingToProcess
	}
	if bs.BatchSize < len(batch) {
		return ErrBlocked
	}
	fmt.Printf("Processed %d elements\n", len(batch))
	<-ctx.Done()
	return nil
}

// Batch is a batch of items.
type Batch []Item

// Item is some abstract item.
type Item struct{}

func main() {
	//initializing BatchService structure
	batchService := BatchService{}
	setFieldToBatchService(InputNumberOfElements, &batchService.NumberOfElements)
	setFieldToBatchService(InputBatchSize, &batchService.BatchSize)
	setFieldToBatchService(InputTimeout, &batchService.Timeout)

	itemSlice := make([]Item, batchService.NumberOfElements)
	for i := 0; i < batchService.NumberOfElements; i++ {
		itemSlice[i] = struct{}{}
	}

	startIndex := 0
	var batch []Item
	for startIndex < batchService.NumberOfElements {
		batch = MakeNextBatch(&itemSlice, startIndex, batchService)
		startIndex += batchService.BatchSize

		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(batchService.Timeout)*time.Second)
		defer cancel()
		err := batchService.Process(timeoutCtx, batch)

		if err == ErrBlocked {
			fmt.Println("Could send batch, service was blocked")
		}
	}
}
