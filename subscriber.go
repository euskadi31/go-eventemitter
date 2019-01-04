// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package eventemitter

// EventSubscriber knows himself what events he is interested in.
// If an EventSubscriber is added to an EventDispatcherInterface,
// the manager invokes {@link SubscribedEvents} and registers
// the subscriber as a listener for all returned events.
type EventSubscriber interface {
	SubscribedEvents() map[string][]interface{}
}
