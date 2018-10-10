package main

import (
	"sync"
	"time"

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
	wg.Add(K + 1)
	go nbuyer.Seller(protocol, K, 1, conn, conn.BasePort(), wg)
	if K == 2 {
		go func(i int) {
			time.Sleep(time.Duration((K-i+1)*100) * time.Millisecond)
			nbuyer.Buyer1_family2(protocol, K, i,
				conn, "localhost", conn.Port(i), // Seller
				conn, "localhost", conn.Port(K+i+1), // ->Buyer[2]
				wg)
		}(1)
		go func(i int) {
			time.Sleep(time.Duration((K-i+1)*100) * time.Millisecond)
			nbuyer.BuyerK_family2(protocol, K, i,
				conn, "localhost", conn.Port(i), // Seller
				conn, conn.Port(K+i), // Buyer[1]->
				wg)
		}(2)
	} else {
		for i := 1; i <= K; i++ {
			if i == 1 {
				go func(i int) {
					time.Sleep(time.Duration((K-i+1)*100) * time.Millisecond)
					nbuyer.Buyer1(protocol, K, i,
						conn, "localhost", conn.Port(i), // Seller
						conn, "localhost", conn.Port(K+i+1), // ->Buyer[i+1]
						wg)
				}(i)
			} else if 1 < i && i < K {
				go func(i int) {
					time.Sleep(time.Duration((K-i+1)*100) * time.Millisecond)
					nbuyer.Buyer(protocol, K, i,
						conn, "localhost", conn.Port(i), // Seller
						conn, conn.Port(K+i), // Buyer[i-1]->
						conn, "localhost", conn.Port(K+i+1), // ->Buyer[i+1]
						wg)
				}(i)
			} else {
				go func(i int) {
					time.Sleep(time.Duration((K-i+1)*100) * time.Millisecond)
					nbuyer.BuyerK(protocol, K, i,
						conn, "localhost", conn.Port(i), // Seller
						conn, conn.Port(K+i), // Buyer[i-1]->
						wg)
				}(i)
			}
		}
	}
	wg.Wait()
}
