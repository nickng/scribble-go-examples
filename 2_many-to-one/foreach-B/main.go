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
	wg.Add(1)
	gather.B(protocol, K, 1, connAB, connAB.BasePort(), wg)
	wg.Wait()
}
