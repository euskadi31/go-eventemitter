// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package eventemitter

import (
	"reflect"
	"sync"
)

// EventEmitter interface
type EventEmitter interface {
	Subscribe(name string, listener interface{})

	Unsubscribe(name string, listener interface{})

	Dispatch(name string, args ...interface{})

	Wait()
}

var _ EventEmitter = (*emitter)(nil)

type emitter struct {
	listeners map[string][]reflect.Value
	sync.RWMutex
	wg sync.WaitGroup
}

// New EventEmitter
func New() EventEmitter {
	return &emitter{
		listeners: make(map[string][]reflect.Value),
		wg:        sync.WaitGroup{},
	}
}

func (e *emitter) Subscribe(name string, listener interface{}) {
	e.Lock()
	defer e.Unlock()

	e.listeners[name] = append(e.listeners[name], reflect.ValueOf(listener))
}

func (e *emitter) Unsubscribe(name string, listener interface{}) {
	e.Lock()
	defer e.Unlock()

	if listeners, ok := e.listeners[name]; ok && len(listeners) > 0 {
		e.removeListener(listeners, name, e.findListenerIdx(listeners, reflect.ValueOf(listener)))
	}
}

func (e *emitter) findListenerIdx(listeners []reflect.Value, listener reflect.Value) int {
	for idx, l := range listeners {
		if l == listener {
			return idx
		}
	}

	return -1
}

func (e *emitter) removeListener(listeners []reflect.Value, name string, idx int) {
	l := len(listeners)

	if !(0 <= idx && idx < l) {
		return
	}

	e.listeners[name] = append(e.listeners[name][:idx], e.listeners[name][idx+1:]...)
}

func (e *emitter) doDispatch(listeners []reflect.Value, arguments []reflect.Value) {
	for _, listener := range listeners {
		listener.Call(arguments)
	}
}

func (e *emitter) buildArguments(args ...interface{}) []reflect.Value {
	arguments := make([]reflect.Value, 0)

	for _, arg := range args {
		arguments = append(arguments, reflect.ValueOf(arg))
	}

	return arguments
}

// Dispatch
func (e *emitter) Dispatch(name string, args ...interface{}) {
	arguments := e.buildArguments(args...)

	e.wg.Add(1)

	go func(arguments []reflect.Value) {
		e.RLock()

		defer e.wg.Done()

		l, ok := e.listeners[name]
		if !ok {
			e.RUnlock()

			return
		}

		listeners := l[:]

		e.RUnlock()

		e.doDispatch(listeners, arguments)
	}(arguments)
}

func (e *emitter) Wait() {
	e.wg.Wait()
}
