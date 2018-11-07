package main

import (
	"sync"

	"github.com/nickng/scribble-go-examples/10_auction/Auction/Protocol"
	"github.com/nickng/scribble-go-examples/10_auction/internal/auction"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	connAB, K := scributil.ParseFlags()
	protocol := Protocol.New()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		auction.Auctioneer(protocol, K, 1, connAB, "localhost", connAB.BasePort(), wg)
	}()
	wg.Wait()
}
