package main

import (
	"sync"
	"flag"

	"github.com/nickng/scribble-go-examples/2_many-to-one/ManyToOne/Gather"
	"github.com/nickng/scribble-go-examples/2_many-to-one/gather"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	var I int
	flag.IntVar(&I, "I", 1, "self ID")

	connAB, K := scributil.ParseFlags()
	protocol := Gather.New()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	gather.A1toK(protocol, K, I, connAB, "localhost", connAB.Port(I), wg)
	wg.Wait()
}
