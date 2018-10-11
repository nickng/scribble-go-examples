package main

import (
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/8_fork-join/ForkJoin/Protocol"
	"github.com/nickng/scribble-go-examples/8_fork-join/partasks"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	connMW, K := scributil.ParseFlags()
	protocol := Protocol.New()

	wg := new(sync.WaitGroup)
	wg.Add(K + 1)
	for k := 1; k <= K; k++ {
		go partasks.Worker(protocol, K, k, connMW, connMW.Port(k), wg)
	}
	go func() {
		time.Sleep(100 * time.Millisecond)
		partasks.Master(protocol, K, 1, connMW, "localhost", connMW.BasePort(), wg)
	}()
	wg.Wait()
}
