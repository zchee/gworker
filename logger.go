// Copyright 2020 The gworker Authors.
// SPDX-License-Identifier: BSD-3-Clause

package gworker

// PanicHandlerFunc is used to handle panics from each worker goroutine.
// if nil, panics will be thrown out again from worker goroutines.
//
// This function handles when follows situation.
//  if p := recover(); p != nil {
//  	if ph := w.pool.options.PanicHandler; ph != nil {
//  		ph(p)
//  	} else {
//  		log.Printf("worker exits from a panic: %v\n", p)
//  		var buf [4096]byte
//  		n := runtime.Stack(buf[:], false)
//  		log.Printf("worker exits from panic: %s\n", string(buf[:n]))
//  	}
// 	}
type PanicHandlerFunc func(p interface{})

// Logger represents a gworker Logger.
type Logger interface {
	Error(msg string, key string, value interface{})
}

// NewPanicHandler return the PanicHandler using zap.Logger.
func NewPanicHandler(logger Logger) PanicHandlerFunc {
	return func(panicErr interface{}) {
		logger.Error("handle panic", "error", panicErr)
	}
}
