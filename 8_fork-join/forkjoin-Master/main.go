package main

import (
	"sync"

	"github.com/nickng/scribble-go-examples/8_fork-join/ForkJoin/Protocol"
	"github.com/nickng/scribble-go-examples/8_fork-join/partasks"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	connMW, K := scributil.ParseFlags()
	protocol := Protocol.New()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	partasks.Master(protocol, K, 1, connMW, "localhost", connMW.BasePort(), wg)
	wg.Wait()
}
