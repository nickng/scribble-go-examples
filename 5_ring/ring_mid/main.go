//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/5_ring/ring_mid
//$ bin/ring_mid.exe -K=4 -I=3

//go:generate scribblec-param.sh ../Ring.scr -d .. -param RingProto github.com/nickng/scribble-go-examples/5_ring/Ring -param-api W

package main

import (
	"encoding/gob"
	"log"
	"sync"
	"flag"

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

	var I int
	flag.IntVar(&I, "I", -1, "self ID (2 <= I <= K-1)")
	 _, K := scributil.ParseFlags() // K >= 3, 1 < I <= K-1

	wg := new(sync.WaitGroup)
	wg.Add(1)

	ring.ServerClient_mid(wg, K, I)

	wg.Wait()
}
