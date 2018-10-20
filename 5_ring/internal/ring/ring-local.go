//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/5_ring/ring-local
//$ bin/ring-local.exe -t=shm

//go:generate scribblec-param.sh ../Ring.scr -d ../ -param RingProto github.com/nickng/scribble-go-examples/5_ring/Ring -param-api W

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
