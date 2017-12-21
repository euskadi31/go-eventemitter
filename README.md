Go Event Emitter ![Last release](https://img.shields.io/github/release/euskadi31/go-eventemitter.svg)
===================

[![Go Report Card](https://goreportcard.com/badge/github.com/euskadi31/go-eventemitter)](https://goreportcard.com/report/github.com/euskadi31/go-eventemitter)

| Branch  | Status | Coverage |
|---------|--------|----------|
| master  | [![Build Status](https://img.shields.io/travis/euskadi31/go-eventemitter/master.svg)](https://travis-ci.org/euskadi31/go-eventemitter) | [![Coveralls](https://img.shields.io/coveralls/euskadi31/go-eventemitter/master.svg)](https://coveralls.io/github/euskadi31/go-eventemitter?branch=master) |
| develop | [![Build Status](https://img.shields.io/travis/euskadi31/go-eventemitter/develop.svg)](https://travis-ci.org/euskadi31/go-eventemitter) | [![Coveralls](https://img.shields.io/coveralls/euskadi31/go-eventemitter/develop.svg)](https://coveralls.io/github/euskadi31/go-eventemitter?branch=develop) |

go-eventemitter is the little and lightweight event emitter library for Go.

Example
-------

```go
package main

import "github.com/euskadi31/go-eventemitter"

func main() {
    emitter := eventemitter.New()

    emitter.Subscribe("test", func() {
        // code
    })

    emitter.Subscribe("count", func(i int) {
        // code
    })

    emitter.Dispatch("test")

    emitter.Dispatch("count", 42)

    emitter.Wait()
}
```

License
-------

go-eventemitter is licensed under [the MIT license](LICENSE.md).
