// Copyright (c) 2024 The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0

package safeexpvar

import (
	"expvar"
	"sync"
)

// mu serializes the check-then-create step below. expvar.Get and
// expvar.NewInt are each safe to call concurrently, but two goroutines can
// still both observe expvar.Get(name) == nil for the same unregistered name
// and both call expvar.NewInt(name); the second call panics ("Reuse of
// exported var name"). The lock is only needed on that slow path, so an
// unlocked fast-path check keeps the common case (name already registered)
// lock-free.
var mu sync.Mutex

func SetInt(name string, value int64) {
	v := expvar.Get(name)
	if v == nil {
		mu.Lock()
		v = expvar.Get(name)
		if v == nil {
			v = expvar.NewInt(name)
		}
		mu.Unlock()
	}
	v.(*expvar.Int).Set(value)
}
