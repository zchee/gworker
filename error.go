// Copyright 2020 The gworker Authors.
// SPDX-License-Identifier: BSD-3-Clause

package gworker

import (
	"errors"

	ants "github.com/panjf2000/ants/v2"
)

// ErrNotImplement retruns the not implement error.
var ErrNotImplement = errors.New("not implement")

// IsInvalidPoolSizeError reports whether the err is ants.ErrInvalidPoolSize which
// returned when setting a negative number as pool capacity.
func IsInvalidPoolSizeError(err error) bool {
	return errors.Is(err, ants.ErrInvalidPoolSize)
}

// IsLackPoolFuncError reports whether the err is ants.ErrLackPoolFunc which
// returned when invokers don't provide function for pool.
func IsLackPoolFuncError(err error) bool {
	return errors.Is(err, ants.ErrLackPoolFunc)
}

// IsInvalidPoolExpiryError reports whether the err is ants.ErrInvalidPoolExpiry which
// returned when setting a negative number as the periodic duration to purge goroutines.
func IsInvalidPoolExpiryError(err error) bool {
	return errors.Is(err, ants.ErrInvalidPoolExpiry)
}

// IsPoolClosedError reports whether the err is ants.ErrPoolClosed which
// returned when submitting task to a closed pool.
func IsPoolClosedError(err error) bool {
	return errors.Is(err, ants.ErrPoolClosed)
}

// IsPoolOverloadError reports whether the err is ants.ErrPoolOverload which
// returned when the pool is full and no workers available.
func IsPoolOverloadError(err error) bool {
	return errors.Is(err, ants.ErrPoolOverload)
}
