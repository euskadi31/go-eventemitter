// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package eventemitter provides tools that allow your application
// components to communicate with each other by dispatching
// events and listening to them.
package eventemitter

import (
	"reflect"
	"sync"
)

// EventEmitter interface.
type EventEmitter interface {
	Subscribe(name string, listener interface{})

	Unsubscribe(name string, listener interface{})

	Dispatch(name string, args ...interface{})

	AddSubscriber(subscriber EventSubscriber)

	RemoveSubscriber(subscriber EventSubscriber)

	Wait()
}

// Emitter struct.
type Emitter struct {
	listeners map[string][]reflect.Value
	sync.RWMutex
	wg sync.WaitGroup
}

// New EventEmitter.
func New() *Emitter {
	return &Emitter{
		listeners: make(map[string][]reflect.Value),
		wg:        sync.WaitGroup{},
	}
}

// Subscribe adds an event listener that listens on the specified events.
func (e *Emitter) Subscribe(name string, listener interface{}) {
	e.Lock()
	defer e.Unlock()

	e.listeners[name] = append(e.listeners[name], reflect.ValueOf(listener))
}

// Unsubscribe removes an event listener from the specified events.
func (e *Emitter) Unsubscribe(name string, listener interface{}) {
	e.Lock()
	defer e.Unlock()

	if listeners, ok := e.listeners[name]; ok && len(listeners) > 0 {
		e.removeListener(listeners, name, e.findListenerIdx(listeners, reflect.ValueOf(listener)))

		if len(e.listeners[name]) == 0 {
			delete(e.listeners, name)
		}
	}
}

// AddSubscriber adds an event subscriber.
func (e *Emitter) AddSubscriber(subscriber EventSubscriber) {
	for name, listeners := range subscriber.SubscribedEvents() {
		for _, listener := range listeners {
			e.Subscribe(name, listener)
		}
	}
}

// RemoveSubscriber removes an event subscriber.
func (e *Emitter) RemoveSubscriber(subscriber EventSubscriber) {
	for name, listeners := range subscriber.SubscribedEvents() {
		for _, listener := range listeners {
			e.Unsubscribe(name, listener)
		}
	}
}

func (e *Emitter) findListenerIdx(listeners []reflect.Value, listener reflect.Value) int {
	for idx, l := range listeners {
		if reflect.DeepEqual(l, listener) {
			return idx
		}
	}

	return -1
}

func (e *Emitter) removeListener(listeners []reflect.Value, name string, idx int) {
	l := len(listeners)

	if idx < 0 || idx >= l {
		return
	}

	e.listeners[name] = append(e.listeners[name][:idx], e.listeners[name][idx+1:]...)
}

func (e *Emitter) doDispatch(listeners []reflect.Value, arguments []reflect.Value) {
	for _, listener := range listeners {
		listener.Call(arguments)
	}
}

func (e *Emitter) buildArguments(args ...interface{}) []reflect.Value {
	arguments := make([]reflect.Value, 0)

	for _, arg := range args {
		arguments = append(arguments, reflect.ValueOf(arg))
	}

	return arguments
}

// Dispatch  an event to all registered listeners.
func (e *Emitter) Dispatch(name string, args ...interface{}) {
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

		listeners := l

		e.RUnlock()

		e.doDispatch(listeners, arguments)
	}(arguments)
}

// Wait all listeners.
func (e *Emitter) Wait() {
	e.wg.Wait()
}
