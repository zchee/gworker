// Copyright 2020 The gworker Authors.
// SPDX-License-Identifier: BSD-3-Clause

package gworker_test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"

	ants "github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"

	"github.com/zchee/gworker"
)

func TestNewPanicHandler(t *testing.T) {
	tests := map[string]struct {
		panicArg interface{}
		wantErr  string
	}{
		"PanicErrorType": {
			panicArg: errors.New("testPanic"),
			wantErr:  `{"error": "testPanic"}`,
		},
		"PanicStructType": {
			panicArg: struct{}{},
			wantErr:  `{"panic": {}}`,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger, buf := NewTestZapLogger(t)
			phfn := gworker.NewPanicHandler(logger)

			var wg sync.WaitGroup
			ph := func(p interface{}) {
				defer wg.Done() // wait for called PanicHandler, this usage testing only.
				phfn(p)
			}
			w, err := gworker.NewWorker(1, ants.WithPanicHandler(ph))
			if err != nil {
				t.Fatalf("faild to initialize NewWorker: %v", err)
			}
			defer w.Release()

			wg.Add(1)
			if err := w.Submit(func() {
				panic(tt.panicArg)
			}); err != nil {
				t.Fatal(err)
			}
			wg.Wait()

			if got := buf.Stripped(); !strings.EqualFold(got, tt.wantErr) {
				t.Fatalf("got %s but want %s", got, tt.wantErr)
			}
		})
	}
}

// Syncer is a spy for the Sync portion of zapcore.WriteSyncer.
type Syncer interface {
	Sync() (err error)
	SetError(err error)
	Called() (called bool)
}

// syncer is a spy for the Sync portion of zapcore.WriteSyncer.
type syncer struct {
	sync.RWMutex

	err    error
	called bool
}

// compile time check whether the syncer implements Syncer interface.
var _ Syncer = (*syncer)(nil)

// Sync records that it was called, then returns the user-supplied error if any.
func (s *syncer) Sync() (err error) {
	s.Lock()
	s.called = true
	err = s.err
	s.Unlock()

	return
}

// SetError sets the error that the Sync method will return.
func (s *syncer) SetError(err error) {
	s.Lock()
	s.err = err
	s.Unlock()
}

// Called reports whether the Sync method was called.
func (s *syncer) Called() (called bool) {
	s.RLock()
	called = s.called
	s.RUnlock()

	return
}

// Buffer is a spy for the Write portion of zapcore.WriteSyncer.
type Buffer interface {
	fmt.Stringer
	io.Writer
}

// LockedBuffer represents a goroutine-safe Buffer.
type LockedBuffer struct {
	Buffer
	Syncer

	mu sync.Mutex
}

// compile time check whether the LockedBuffer implements zapcore.WriteSyncer interface.
var _ zapcore.WriteSyncer = (*LockedBuffer)(nil)

// LockedBufferOption represents the optional function.
type LockedBufferOption func(lb *LockedBuffer)

// WithBuffer inject Buffer to LockedBuffer.
func WithBuffer(b Buffer) LockedBufferOption {
	return func(lb *LockedBuffer) {
		lb.Buffer = b
	}
}

// WithSyncer inject Syncer to LockedBuffer.
func WithSyncer(s Syncer) LockedBufferOption {
	return func(lb *LockedBuffer) {
		lb.Syncer = s
	}
}

var (
	_ = WithBuffer
	_ = WithSyncer
)

// NewLockedBuffer returns the new LockedBuffer.
func NewLockedBuffer(opts ...LockedBufferOption) *LockedBuffer {
	lb := &LockedBuffer{
		Buffer: new(strings.Builder),
		Syncer: new(syncer),
	}

	for _, o := range opts {
		o(lb)
	}

	return lb
}

// Write implements zapcore.WriteSyncer.
func (lb *LockedBuffer) Write(bs []byte) (n int, err error) {
	lb.mu.Lock()
	n, err = lb.Buffer.Write(bs)
	lb.mu.Unlock()
	return
}

// Sync implements zapcore.WriteSyncer.
func (lb *LockedBuffer) Sync() (err error) {
	lb.mu.Lock()
	err = lb.Syncer.Sync()
	lb.mu.Unlock()
	return
}

// Lines returns the current buffer contents, split on newlines.
//
// mimic to zap/ztest.Buffer with goroutine-safe.
func (lb *LockedBuffer) Lines() (output []string) {
	lb.mu.Lock()
	output = strings.Split(lb.Buffer.String(), "\n")
	output = output[:len(output)-1]
	lb.mu.Unlock()

	return
}

// Stripped returns the current buffer contents with the last trailing newline
// stripped.
//
// mimic to zap/ztest.Buffer with goroutine-safe.
func (lb *LockedBuffer) Stripped() (s string) {
	lb.mu.Lock()
	s = strings.TrimRight(lb.Buffer.String(), "\n")
	lb.mu.Unlock()

	return
}

// NewTestZapLogger returns the new zaptest.NewLogger and LockedBuffer.
func NewTestZapLogger(t *testing.T, lbopts ...LockedBufferOption) (*zap.Logger, *LockedBuffer) {
	buf := NewLockedBuffer(lbopts...)
	opts := []zaptest.LoggerOption{zaptest.WrapOptions(
		zap.WrapCore(func(zapcore.Core) zapcore.Core {
			return zapcore.NewCore(
				zapcore.NewConsoleEncoder(zapcore.EncoderConfig{}), // output only entries
				buf, // avoid occur race data
				zapcore.ErrorLevel,
			)
		}),
	)}

	logger := zaptest.NewLogger(t, opts...)
	t.Cleanup(func() {
		if err := logger.Sync(); err != nil {
			t.Fatal(err)
		}
	})

	return logger, buf
}
