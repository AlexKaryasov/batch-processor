package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrBlocked reports if service is blocked.
var ErrBlocked = errors.New("blocked")

// Service defines external service that can process batches of items.
type Service interface {
	Process(ctx context.Context, batch Batch) error
}

type BatchService struct {
	NumberOfElements int
	BatchSize        int
	Timeout          int
}

func (bs BatchService) Process(ctx context.Context, batch Batch) error {
	if bs.BatchSize < len(batch) {
		return ErrBlocked
	}
	fmt.Printf("Processing batch with size %d within %d seconds\n", bs.BatchSize, bs.Timeout)
	fmt.Printf("Processed %d elements\n", len(batch))
	<-ctx.Done()
	return nil
}

// Batch is a batch of items.
type Batch []Item

// Item is some abstract item.
type Item struct{}

func main() {
	batchService := BatchService{}
	for batchService.NumberOfElements == 0 {
		fmt.Println("Please input how many elements are you going to process")
		fmt.Scan(&batchService.NumberOfElements)
	}
	for batchService.BatchSize == 0 {
		fmt.Println("Please input batch size")
		_, _ = fmt.Scan(&batchService.BatchSize)
	}
	for batchService.Timeout == 0 {
		fmt.Println("Please input batch timeout")
		fmt.Scan(&batchService.Timeout)
	}

	itemSlice := make([]Item, batchService.NumberOfElements)
	for i := 0; i < batchService.NumberOfElements; i++ {
		itemSlice[i] = struct{}{}
	}

	startIndex := 0
	var batch []Item
	for startIndex < batchService.NumberOfElements {
		batch = makeNextBatch(&itemSlice, startIndex, batchService)
		timeoutCtx, _ := context.WithTimeout(context.Background(), time.Duration(batchService.Timeout)*time.Second)
		startIndex += batchService.BatchSize
		err := batchService.Process(timeoutCtx, batch)
		if err == ErrBlocked {
			fmt.Println("Could send batch, service was blocked")
		}
	}
}

func makeNextBatch(items *[]Item, startIndex int, batchService BatchService) []Item {
	if startIndex+batchService.BatchSize >= batchService.NumberOfElements {
		return (*items)[startIndex:]
	}
	return (*items)[startIndex : startIndex+batchService.BatchSize]
}
