//rhu@HZHL4 ~/code/go
//$ go install github.com/nickng/scribble-go-examples/1_one-to-many/scatter-local
//$ bin/scatter-local.exe

//go:generate scribblec-param.sh ../OneToMany.scr -d ../ -param Scatter github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany -param-api A -param-api B


package main

import (
	"encoding/gob"
	"sync"
	"time"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"

	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter"
	"github.com/nickng/scribble-go-examples/1_one-to-many/messages"
	"github.com/nickng/scribble-go-examples/1_one-to-many/scatter"
)

var _ = shm.Listen
var _ = tcp.Listen

func init() {
	var data onetomany.Data
	gob.Register(&data)
}

func main() {
	port := 33333
	K := 2
	p := Scatter.New()  // FIXME: K should be param here?
	wg := new(sync.WaitGroup)
	wg.Add(2)
	/*/
	listen := tcp.BListen
	dial := tcp.Dial
	fmtr := func() session2.ScribMessageFormatter { return new(session2.GobFormatter) } 
	/*/
	listen := shm.BListen
	dial := shm.Dial
	fmtr := func() session2.ScribMessageFormatter { return new(session2.PassByPointer) }
	//*/
	for i := 1; i <= K; i++ {
		go scatter.Server_gather(listen, fmtr, port+i, p, K, i, wg)
	}
	time.Sleep(100 * time.Millisecond) //2017/12/11 11:21:40 cannot connect to 127.0.0.1:8891: dial tcp 127.0.0.1:8891: connectex: No connection could be made because the target machine actively refused it.
	scatter.Client_scatter(dial, fmtr, "localhost", port, p, K, 1)
	wg.Wait()
}
