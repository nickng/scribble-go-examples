package main

import (
	"sync"

	"github.com/nickng/scribble-go-examples/21_nbuyers/NBuyers/NBuyers"
	nbuyer "github.com/nickng/scribble-go-examples/21_nbuyers/internal/nbuyers"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	conn, K := scributil.ParseFlags()
	protocol := NBuyers.New()

	// Port offsets:
	// 1..K: Buyer{1,2..K-1,K} → Seller
	// K+2 : Buyer[1] → Buyer[2]
	// K+i+1 : Buyer[i] → Buyer[i+1]
	// K+K : Buyer[K-1] → Buyer[K]

	wg := new(sync.WaitGroup)
	wg.Add(1)
	nbuyer.Seller(protocol, K, 1, conn, conn.BasePort(), wg)
	wg.Wait()
}
