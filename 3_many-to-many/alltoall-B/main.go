package main

import (
	"sync"
	"flag"

	"github.com/nickng/scribble-go-examples/3_many-to-many/ManyToMany/Alltoall"
	"github.com/nickng/scribble-go-examples/3_many-to-many/alltoall"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	var I int
	flag.IntVar(&I, "I", 1, "self ID") // 1<= I <= M
	var M int
	flag.IntVar(&M, "M", 1, "num of 'A' workers") // M + N = K
	var N int
	flag.IntVar(&N, "N", 1, "num of 'B' workers") // M + N = K

	connAB, K := scributil.ParseFlags()
	// M, N := alltoall.SplitMN(K)
	scributil.Debugf("[info] K=%d: M=%d N=%d\n", K, M, N)
	protocol := Alltoall.New()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	alltoall.B(protocol, M, N, I, connAB, connAB.BasePort(), wg)
	wg.Wait()
}
