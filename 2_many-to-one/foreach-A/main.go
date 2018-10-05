package main

import (
	"sync"

	"github.com/nickng/scribble-go-examples/2_many-to-one/ManyToOne/Gather"
	"github.com/nickng/scribble-go-examples/2_many-to-one/gather"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	connAB, K := scributil.ParseFlags()
	protocol := Gather.New()

	wg := new(sync.WaitGroup)
	wg.Add(K)
	for i := 1; i <= K; i++ {
		go gather.A1toK(protocol, K, i, connAB, "localhost", connAB.Port(i), wg)
	}
	wg.Wait()
}
