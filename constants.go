// Copyright 2020 The gworker Authors.
// SPDX-License-Identifier: BSD-3-Clause

package gworker

import (
	ants "github.com/panjf2000/ants/v2"
)

const (
	// DefaultAntsPoolSize is the default capacity for a default goroutine pool.
	//
	// actual size is math.MaxInt32.
	DefaultAntsPoolSize = ants.DefaultAntsPoolSize

	// DefaultCleanIntervalTime is the interval time to clean up goroutines.
	//
	// actual interval time is time.Second.
	DefaultCleanIntervalTime = ants.DefaultCleanIntervalTime
)
