// Copyright 2020 The gworker Authors.
// SPDX-License-Identifier: BSD-3-Clause

package gworker_test

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zchee/gworker"
)

func ExampleNewWorker() {
	testSubmitFunc := func(args *int64) {
		atomic.AddInt64(args, 1)
	}

	w, err := gworker.NewWorker(1000)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Release()

	const loopCount = 100
	var wg sync.WaitGroup
	arg := int64(0) // atmoic
	for i := 0; i < loopCount; i++ {
		wg.Add(1)
		if err := w.Submit(func() {
			defer wg.Done()
			testSubmitFunc(&arg)
		}); err != nil {
			log.Fatal(err)
		}
	}
	wg.Wait()

	fmt.Println(atomic.LoadInt64(&arg) == loopCount) // increase args per loopCount atomically

	// Output:
	// true
}

type Bill struct {
	UserID    int64     `json:"UserId"`
	Status    string    `json:"Status"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type invoker struct {
	wg *sync.WaitGroup
}

func (invoker) getBill(ctx context.Context) *Bill {
	return &Bill{
		UserID:    0,
		Status:    "buy",
		CreatedAt: time.Unix(0, 0).UTC(),
	}
}

func (i *invoker) work(iface interface{}) {
	defer i.wg.Done()

	bill := iface.(*Bill)
	fmt.Printf("UserID: %d, Status: %s, CreatedAt: %s\n", bill.UserID, bill.Status, bill.CreatedAt)
}

func ExampleNewWorkerFunc() {
	ctx := context.Background()

	var wg sync.WaitGroup
	invoker := &invoker{&wg}

	w, err := gworker.NewWorkerFunc(1000, invoker.work)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Release()

	wg.Add(1)
	if err := w.Invoke(invoker.getBill(ctx)); err != nil {
		log.Fatal(err)
	}
	wg.Wait()

	// Output:
	// UserID: 0, Status: buy, CreatedAt: 1970-01-01 00:00:00 +0000 UTC
}
