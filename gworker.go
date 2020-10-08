// Copyright 2020 The gworker Authors.
// SPDX-License-Identifier: BSD-3-Clause

package gworker

import (
	ants "github.com/panjf2000/ants/v2"
)

// Worker represents a goroutine worker.
type Worker interface {
	// Running returns the number of the currently running goroutines.
	Running() int

	// Free returns a available goroutines to work.
	Free() int

	// Cap returns the capacity of this pool.
	Cap() int

	// Tune changes the capacity of this pool.
	Tune(size int)

	// Release closes this pool.
	Release()

	// Reboot reboots a released pool.
	Reboot()

	// Submit submits a task to this pool.
	Submit(task func()) error

	// Invoke submits a task to pool.
	Invoke(args interface{}) error
}

// worker represents a Submit style goroutine worker.
type worker struct {
	pool *ants.Pool
}

// compile time check whether the worker implements Worker interface.
var _ Worker = (*worker)(nil)

// NewWorker returns the new Worker which Submit style goroutine worker.
func NewWorker(size int32, options ...ants.Option) (Worker, error) {
	p, err := ants.NewPool(int(size), options...)
	if err != nil {
		return nil, err
	}

	return &worker{pool: p}, nil
}

// Running implements Worker.
func (w *worker) Running() int { return w.pool.Running() }

// Free implements Worker.
func (w *worker) Free() int { return w.pool.Free() }

// Cap implements Worker.
func (w *worker) Cap() int { return w.pool.Cap() }

// Tune implements Worker.
func (w *worker) Tune(size int) { w.pool.Tune(size) }

// Release implements Worker.
func (w *worker) Release() { w.pool.Release() }

// Reboot implements Worker.
func (w *worker) Reboot() { w.pool.Reboot() }

// Submit implements Worker.
func (w *worker) Submit(task func()) error { return w.pool.Submit(task) }

// Invoke implements Worker.
//
// Invoke not implemented worker, will returns the ErrNotImplement error.
func (w *worker) Invoke(args interface{}) error { return ErrNotImplement }

// workerFunc represents a Invoke style goroutine worker.
type workerFunc struct {
	pool *ants.PoolWithFunc
}

// compile time check whether the workerFunc implements Worker interface.
var _ Worker = (*workerFunc)(nil)

// Func represents a worker function.
type Func func(args interface{})

// NewWorkerFunc returns the new Worker which Invoke style goroutine worker.
func NewWorkerFunc(size int32, fn Func, options ...ants.Option) (Worker, error) {
	p, err := ants.NewPoolWithFunc(int(size), fn, options...)
	if err != nil {
		return nil, err
	}

	return &workerFunc{pool: p}, nil
}

// Running implements Worker.
func (w *workerFunc) Running() int { return w.pool.Running() }

// Free implements Worker.
func (w *workerFunc) Free() int { return w.pool.Free() }

// Cap implements Worker.
func (w *workerFunc) Cap() int { return w.pool.Cap() }

// Tune implements Worker.
func (w *workerFunc) Tune(size int) { w.pool.Tune(size) }

// Release implements Worker.
func (w *workerFunc) Release() { w.pool.Release() }

// Reboot implements Worker.
func (w *workerFunc) Reboot() { w.pool.Reboot() }

// Submit implements Worker.
//
// Submit not implemented workerFunc, will returns the ErrNotImplement error.
func (w *workerFunc) Submit(task func()) error { return ErrNotImplement }

// Invoke implements Worker.
func (w *workerFunc) Invoke(args interface{}) error { return w.pool.Invoke(args) }
