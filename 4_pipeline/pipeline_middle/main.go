package main

import (
	"encoding/gob"
	"log"
	"sync"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"

	"github.com/nickng/scribble-go-examples/4_pipeline/messages"
	"github.com/nickng/scribble-go-examples/4_pipeline/pipeline"
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
}



func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	K := 4

	wg := new(sync.WaitGroup)
	wg.Add(K-2)

	for j := 2; j <= K-1; j++ {
		go pipeline.Server_middle(wg, K, j)
	}

	wg.Wait()
}
