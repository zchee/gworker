// Copyright 2020 The gworker Authors.
// SPDX-License-Identifier: BSD-3-Clause

package gworker_test

import (
	"errors"
	"testing"

	ants "github.com/panjf2000/ants/v2"

	"github.com/zchee/gworker"
)

func TestIsError(t *testing.T) {
	t.Parallel()

	var notAntsError = errors.New("testNotAntsError")

	tests := map[string]struct {
		errFn func(err error) bool
		err   error
		want  bool
	}{
		"HandleAntsInvalidPoolSizeError": {
			errFn: gworker.IsInvalidPoolSizeError,
			err:   ants.ErrInvalidPoolSize,
			want:  true,
		},
		"HandleNotAntsInvalidPoolSizeError": {
			errFn: gworker.IsInvalidPoolSizeError,
			err:   notAntsError,
			want:  false,
		},
		"HandleAntsLackPoolFuncError": {
			errFn: gworker.IsLackPoolFuncError,
			err:   ants.ErrLackPoolFunc,
			want:  true,
		},
		"HandleNotAntsLackPoolFuncError": {
			errFn: gworker.IsLackPoolFuncError,
			err:   notAntsError,
			want:  false,
		},
		"HandleAntsInvalidPoolExpiryError": {
			errFn: gworker.IsInvalidPoolExpiryError,
			err:   ants.ErrInvalidPoolExpiry,
			want:  true,
		},
		"HandleNotAntsInvalidPoolExpiryError": {
			errFn: gworker.IsInvalidPoolExpiryError,
			err:   notAntsError,
			want:  false,
		},
		"HandleAntsPoolClosed": {
			errFn: gworker.IsPoolClosedError,
			err:   ants.ErrPoolClosed,
			want:  true,
		},
		"HandleNotAntsPoolClosed": {
			errFn: gworker.IsPoolClosedError,
			err:   notAntsError,
			want:  false,
		},
		"HandleAntsPoolOverloadError": {
			errFn: gworker.IsPoolOverloadError,
			err:   ants.ErrPoolOverload,
			want:  true,
		},
		"HandleNotAntsPoolOverloadError": {
			errFn: gworker.IsPoolOverloadError,
			err:   notAntsError,
			want:  false,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := tt.errFn(tt.err); got != tt.want {
				t.Errorf("got %v but want %v error: %#v", got, tt.want, tt.err)
			}
		})
	}
}
