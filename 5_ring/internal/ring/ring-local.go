package ring

import (
	"sync"
	"time"

	//"github.com/nickng/scribble-go-examples/5_ring/Ring/RingProto"
	"github.com/nickng/scribble-go-examples/scributil"
)

func Local() {
	_, K := scributil.ParseFlags()

	//p := RingProto.New()
	wg := new(sync.WaitGroup)
	wg.Add(K)
	go Server_last(wg, K, K)
	for i := 2; i < K; i++ {
		go ServerClient_mid(wg, K, i)
	}
	time.Sleep(100 * time.Millisecond)
	Client_ini(wg, K, 1)
	wg.Wait()
}
