package main

import (
	"sync"
	"flag"

	"github.com/nickng/scribble-go-examples/8_fork-join/ForkJoin/Protocol"
	"github.com/nickng/scribble-go-examples/8_fork-join/partasks"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	var I int
	flag.IntVar(&I, "I", -1, "I>2")
	connMW, K := scributil.ParseFlags()
	protocol := Protocol.New()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	partasks.Worker(protocol, K, I, connMW, connMW.Port(I), wg)
	wg.Wait()
}
