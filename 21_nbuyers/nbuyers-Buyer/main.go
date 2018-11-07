package main

import (
	"flag"
	"sync"

	"github.com/nickng/scribble-go-examples/21_nbuyers/NBuyers/NBuyers"
	nbuyer "github.com/nickng/scribble-go-examples/21_nbuyers/internal/nbuyers"
	"github.com/nickng/scribble-go-examples/scributil"
)

var buyer = flag.Int("I", 1, "Specify Buyer ID (1..K)")

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
	if K == 2 {
		i := *buyer
		if *buyer == 1 {
			i := *buyer
			nbuyer.Buyer1_family2(protocol, K, i,
				conn, "localhost", conn.Port(i), // Seller
				conn, "localhost", conn.Port(K+i+1), // ->Buyer[2]
				wg)
		} else {
			nbuyer.BuyerK_family2(protocol, K, i,
				conn, "localhost", conn.Port(i), // Seller
				conn, conn.Port(K+i), // Buyer[1]->
				wg)
		}
	} else { // K > 2
		i := *buyer
		if i == 1 {
			nbuyer.Buyer1(protocol, K, i,
				conn, "localhost", conn.Port(i), // Seller
				conn, "localhost", conn.Port(K+i+1), // ->Buyer[i+1]
				wg)
		} else if 1 < i && i < K {
			nbuyer.Buyer(protocol, K, i,
				conn, "localhost", conn.Port(i), // Seller
				conn, conn.Port(K+i), // Buyer[i-1]->
				conn, "localhost", conn.Port(K+i+1), // ->Buyer[i+1]
				wg)
		} else {
			nbuyer.BuyerK(protocol, K, i,
				conn, "localhost", conn.Port(i), // Seller
				conn, conn.Port(K+i), // Buyer[i-1]->
				wg)
		}
	}
	wg.Wait()
}
