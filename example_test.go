// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package eventemitter_test

import (
	"fmt"

	"github.com/euskadi31/go-eventemitter"
)

func ExampleEmitter_Subscribe() {
	emitter := eventemitter.New()

	// Subscribe user.created event
	emitter.Subscribe("user.created", func(id int) {
		fmt.Printf("UserID: %d\n", id)
	})

	// Dispatch user.created event with user id
	emitter.Dispatch("user.created", 1234)

	// wait for all listeners to be executed
	emitter.Wait()
}

type userSubscriber struct{}

func (s userSubscriber) SubscribedEvents() map[string][]interface{} {
	return map[string][]interface{}{
		"user.created": {
			s.onCreateUser(),
		},
	}
}

func (s *userSubscriber) onCreateUser() func(id int) {
	return func(id int) {
		fmt.Printf("UserID: %d\n", id)
	}
}

func ExampleEmitter_AddSubscriber() {
	//see: example_test.go#L30 for userSubscriber

	emitter := eventemitter.New()

	// Subscribe user.created event
	emitter.AddSubscriber(&userSubscriber{})

	// Dispatch user.created event with user id
	emitter.Dispatch("user.created", 1234)

	// wait for all listeners to be executed
	emitter.Wait()
}
