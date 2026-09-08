// Copyright (c) 2024 The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0

package safeexpvar

import (
	"expvar"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jaegertracing/jaeger/internal/testutils"
)

func TestSetInt(t *testing.T) {
	// Test with a new variable
	name := "metrics-test-1"
	value := int64(42)

	SetInt(name, value)

	// Retrieve the variable and check its value
	v := expvar.Get(name)
	assert.NotNil(t, v, "expected variable %s to be created", name)
	expInt, ok := v.(*expvar.Int)
	require.True(t, ok, "expected variable %s to be of type *expvar.Int", name)
	assert.Equal(t, value, expInt.Value())
}

// TestSetIntConcurrentFirstRegistration reproduces the panic reported in
// #9282: many goroutines racing SetInt on the SAME not-yet-registered name
// used to both observe expvar.Get(name) == nil and both call
// expvar.NewInt(name), and the second call panics ("Reuse of exported var
// name"). Run with -race to also confirm no data race on the returned
// *expvar.Int.
func TestSetIntConcurrentFirstRegistration(t *testing.T) {
	for i := range 20 {
		name := fmt.Sprintf("metrics-test-concurrent-%d", i)
		var wg sync.WaitGroup
		for g := range 50 {
			value := int64(g)
			wg.Go(func() {
				SetInt(name, value)
			})
		}
		wg.Wait()

		v := expvar.Get(name)
		require.NotNil(t, v, "expected variable %s to be created", name)
		_, ok := v.(*expvar.Int)
		assert.True(t, ok, "expected variable %s to be of type *expvar.Int", name)
	}
}

func TestMain(m *testing.M) {
	testutils.VerifyGoLeaks(m)
}
