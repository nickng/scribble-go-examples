//go:generate scribblec-param.sh ../Ring.scr -d .. -param RingProto github.com/nickng/scribble-go-examples/5_ring/Ring -param-api W

package main

import (
	"encoding/gob"
	"log"
	"sync"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"

	"github.com/nickng/scribble-go-examples/5_ring/internal/ring"
	"github.com/nickng/scribble-go-examples/5_ring/messages"
	"github.com/nickng/scribble-go-examples/scributil"
)

var _ = shm.Dial
var _ = tcp.Dial

//*
var LISTEN = tcp.Listen
var DIAL = tcp.Dial
var FORMATTER = func() *session2.GobFormatter { return new(session2.GobFormatter) }

/*/
var LISTEN = shm.Listen
var DIAL = shm.Dial
var FORMATTER = func() *session2.PassByPointer { return new(session2.PassByPointer) }
//*/

const PORT = 8888

func init() {
	var foo messages.Foo
	gob.Register(&foo)
	var bar messages.Bar
	gob.Register(&bar)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	_, K := scributil.ParseFlags() // K >= 3

	wg := new(sync.WaitGroup)
	wg.Add(1)

	ring.Ring_last(wg, K, K)

	wg.Wait()
}
