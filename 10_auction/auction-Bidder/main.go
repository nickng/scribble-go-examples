package main

import (
	"sync"
	"flag"

	"github.com/nickng/scribble-go-examples/10_auction/Auction/Protocol"
	"github.com/nickng/scribble-go-examples/10_auction/internal/auction"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	var I int
	flag.IntVar(&I, "I", 1, "self ID")
	connAB, K := scributil.ParseFlags()
	protocol := Protocol.New()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	auction.Bidder(protocol, K, I, connAB, connAB.Port(I), wg)
	wg.Wait()
}
