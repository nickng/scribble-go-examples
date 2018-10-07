package main

import (
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/10_auction/Auction/Protocol"
	"github.com/nickng/scribble-go-examples/10_auction/internal/auction"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	connAB, K := scributil.ParseFlags()
	protocol := Protocol.New()

	wg := new(sync.WaitGroup)
	wg.Add(K + 1)
	for k := 1; k <= K; k++ {
		go auction.Bidder(protocol, K, k, connAB, connAB.Port(k), wg)
	}
	go func() {
		time.Sleep(100 * time.Millisecond)
		auction.Auctioneer(protocol, K, 1, connAB, "localhost", connAB.BasePort(), wg)
	}()
	wg.Wait()
}
