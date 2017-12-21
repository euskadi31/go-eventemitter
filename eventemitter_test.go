// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package eventemitter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmitter(t *testing.T) {
	called := 0

	em := New()

	em.Dispatch("test-empty")

	em.Subscribe("test", func() {
		called++
	})

	em.Subscribe("test", func() {
		called++
	})

	listener := func() {
		called++
	}

	em.Subscribe("test", listener)

	em.Unsubscribe("test", listener)

	em.Dispatch("test")

	em.Wait()

	assert.Equal(t, 2, called)
}

func TestDispatchWithArguments(t *testing.T) {
	called := 0

	em := New()

	em.Subscribe("test", func(i int) {
		called += i
	})

	em.Subscribe("test", func(i int) {
		called += i
	})

	listener := func(i int) {
		called += i
	}

	em.Subscribe("test", listener)

	em.Unsubscribe("test", listener)

	em.Dispatch("test", 10)

	em.Wait()

	assert.Equal(t, 20, called)
}

func TestUnsubscribe(t *testing.T) {
	called := 0

	em := New()

	em.Subscribe("test", func() {
		called++
	})

	listener := func() {
		called++
	}

	em.Unsubscribe("test", listener)

	em.Unsubscribe("test-empty", listener)

	em.Dispatch("test")

	em.Wait()

	assert.Equal(t, 1, called)
}

func BenchmarkDispatchParallel(b *testing.B) {
	em := New()

	em.Subscribe("test-1", func() {})
	em.Subscribe("test-1", func() {})

	em.Subscribe("test-2", func() {})
	em.Subscribe("test-2", func() {})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			em.Dispatch("test-1")
		}
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			em.Dispatch("test-2")
		}
	})

	em.Wait()
}
