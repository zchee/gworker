// Copyright 2020 The gworker Authors.
// SPDX-License-Identifier: BSD-3-Clause

package gworker_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/zchee/gworker"
)

const (
	testWorkerSize = 100
	testNumWorker  = 100
)

func testSubmitFunc(tb testing.TB, args *int64) {
	atomic.AddInt64(args, 1)
}

func TestNewWorker(t *testing.T) {
	w, err := gworker.NewWorker(testWorkerSize)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Release()

	wg := new(sync.WaitGroup)
	errc := make(chan error, testNumWorker) // use buffered chan
	args := int64(0)                        // atomic
	for i := 0; i < testNumWorker; i++ {
		wg.Add(1)
		errc <- w.Submit(func() {
			testSubmitFunc(t, &args)
			wg.Done()
		})
	}

	go func() {
		// wait blocks until the WaitGroup counter is zero and
		// close errc for starting for range errc error check.
		wg.Wait()
		close(errc)
	}()

	// starts for range loop after the closed errc chan
	for err := range errc {
		if err != nil {
			t.Fatal(err)
		}
	}

	if got := atomic.LoadInt64(&args); got != testNumWorker { // increase args per testNumWorker atomically
		t.Fatalf("got %d but want %d", got, testNumWorker)
	}
}

type testInvoke struct {
	mu sync.Mutex
	wg *sync.WaitGroup
}

func (i *testInvoke) work(args interface{}) {
	i.mu.Lock()
	defer i.mu.Unlock()

	n := args.(*int64)
	*n++
	i.wg.Done()
}

func TestNewWorkerFunc(t *testing.T) {
	wg := new(sync.WaitGroup)
	invoker := &testInvoke{wg: wg}
	w, err := gworker.NewWorkerFunc(testWorkerSize, invoker.work)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Release()

	errc := make(chan error, testNumWorker) // use buffered chan
	args := int64(0)                        // atomic
	for i := 0; i < testNumWorker; i++ {
		wg.Add(1)
		errc <- w.Invoke(&args)
	}

	go func() {
		// wait blocks until the WaitGroup counter is zero and
		// close errc for starting for range errc error check.
		wg.Wait()
		close(errc)
	}()

	// starts for range loop after the closed errc chan
	for err := range errc {
		if err != nil {
			t.Fatal(err)
		}
	}

	if got := atomic.LoadInt64(&args); got != testNumWorker { // increase args per testNumWorker atomically
		t.Fatalf("got %d but want %d", got, testNumWorker)
	}
}

func TestWorkerCheckSize(t *testing.T) {
	tests := map[string]struct {
		numWorker   int32
		wantRunning int
		wantCap     int
		wantFree    int
		tuned       bool
		tuneSize    int
	}{
		"RunOneWorker": {
			numWorker:   1,
			wantRunning: 1,
			wantCap:     1,
			wantFree:    0,
		},
		"RunOneWorkerWithTune": {
			numWorker:   1,
			wantRunning: 1,
			wantCap:     5,
			wantFree:    4,
			tuned:       true,
			tuneSize:    5,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w, err := gworker.NewWorker(tt.numWorker)
			if err != nil {
				t.Fatal(err)
			}
			defer w.Release()

			wg := new(sync.WaitGroup)
			wg.Add(1)
			_ = w.Submit(func() {
				testSubmitFunc(t, new(int64))
				wg.Done()
			})

			if tt.tuned {
				w.Tune(tt.tuneSize)
			}
			if gotRunning := w.Running(); gotRunning != tt.wantRunning {
				t.Fatalf("Running: got %d but want %d", gotRunning, tt.wantRunning)
			}
			if gotCap := w.Cap(); gotCap != tt.wantCap {
				t.Fatalf("Cap: got %d but want %d", gotCap, tt.wantCap)
			}
			if gotFree := w.Free(); gotFree != tt.wantFree {
				t.Fatalf("Free: got %d but want %d", gotFree, tt.wantFree)
			}

			wg.Wait()
		})
	}
}

func TestWorkerFuncCheckSize(t *testing.T) {
	tests := map[string]struct {
		numWorker   int32
		wantRunning int
		wantCap     int
		wantFree    int
		tuned       bool
		tuneSize    int
	}{
		"RunOneWorker": {
			numWorker:   1,
			wantRunning: 1,
			wantCap:     1,
			wantFree:    0,
		},
		"RunOneWorkerWithTune": {
			numWorker:   1,
			wantRunning: 1,
			wantCap:     5,
			wantFree:    4,
			tuned:       true,
			tuneSize:    5,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			wg := new(sync.WaitGroup)
			invoker := &testInvoke{wg: wg}
			w, err := gworker.NewWorkerFunc(tt.numWorker, invoker.work)
			if err != nil {
				t.Fatal(err)
			}
			defer w.Release()

			wg.Add(1)
			_ = w.Invoke(new(int64))

			if tt.tuned {
				w.Tune(tt.tuneSize)
			}
			if gotRunning := w.Running(); gotRunning != tt.wantRunning {
				t.Fatalf("Running: got %d but want %d", gotRunning, tt.wantRunning)
			}
			if gotCap := w.Cap(); gotCap != tt.wantCap {
				t.Fatalf("Cap: got %d but want %d", gotCap, tt.wantCap)
			}
			if gotFree := w.Free(); gotFree != tt.wantFree {
				t.Fatalf("Free: got %d but want %d", gotFree, tt.wantFree)
			}

			wg.Wait()
		})
	}
}

func TestWorkerNotImplement(t *testing.T) {
	const numWorker = 1
	MostWorker := func(size int32) gworker.Worker {
		w, err := gworker.NewWorker(size)
		if err != nil {
			panic(err)
		}
		return w
	}

	tests := map[string]struct {
		worker  gworker.Worker
		wantErr error
	}{
		"WorkerWithRunInvoke": {
			worker:  MostWorker(numWorker),
			wantErr: gworker.ErrNotImplement,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := tt.worker
			defer w.Release()

			if err := w.Invoke(nil); !errors.Is(err, tt.wantErr) {
				t.Fatalf("Invoke: got %v but want %v", err, tt.wantErr)
			}
		})
	}
}

func TestWorkerFuncNotImplement(t *testing.T) {
	const numWorker = 1
	MostWorkerFunc := func(size int32) gworker.Worker {
		invoker := &testInvoke{wg: new(sync.WaitGroup)}
		w, err := gworker.NewWorkerFunc(size, invoker.work)
		if err != nil {
			panic(err)
		}
		return w
	}

	tests := map[string]struct {
		worker  gworker.Worker
		wantErr error
	}{
		"WorkerFuncWithRunSummit": {
			worker:  MostWorkerFunc(numWorker),
			wantErr: gworker.ErrNotImplement,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := tt.worker
			defer w.Release()

			if err := w.Submit(nil); !errors.Is(err, tt.wantErr) {
				t.Fatalf("Invoke: got %v but want %v", err, tt.wantErr)
			}
		})
	}
}
