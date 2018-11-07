package main

import (
	"sync"

	"github.com/nickng/scribble-go-examples/2_many-to-one/ManyToOne/Gather"
	"github.com/nickng/scribble-go-examples/2_many-to-one/gather"
	"github.com/nickng/scribble-go-examples/scributil"
)

// Shared memory implementation of ManyToOne.Gather protocol.

func main() {
	connAB, K := scributil.ParseFlags()
	protocol := Gather.New()

	gather.WaitConn = make(chan struct{}) // Wait for servers to be ready
	wg := new(sync.WaitGroup)
	wg.Add(K + 1)
	for i := 1; i <= K; i++ {
		go gather.A1toK(protocol, K, i, connAB, "localhost", connAB.Port(i), wg)
	}
	go gather.B(protocol, K, 1, connAB, connAB.BasePort(), wg)
	wg.Wait()
}
