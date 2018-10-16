package main

import (
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/19_game/Game/Game"
	"github.com/nickng/scribble-go-examples/19_game/Game/Proto1"
	"github.com/nickng/scribble-go-examples/19_game/internal/game"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	conn, K := scributil.ParseFlags()
	proto1 := Proto1.New()
	gameProto := Game.New()

	scributil.Debugf("[info] This example can only work in shm.\n")

	// Port offsets:
	// 1..K        : A -> B (* K)
	// K+1..K+K    : A -> C (* K)
	// 2K+1.. 2K+K : Q -> P[1..K]

	wg := new(sync.WaitGroup)
	wg.Add(2*K + K + 1)
	for i := 1; i <= K; i++ {
		go game.B(gameProto, K, 1, conn, conn.Port(i), wg)
		go game.C(gameProto, K, 1, conn, conn.Port(K+i), wg)
	}
	time.Sleep(100 * time.Millisecond)
	for i := 1; i <= K; i++ {
		go game.P1K(proto1, K, i, conn, conn.Port(2*K+i), wg)
	}
	time.Sleep(100 * time.Millisecond)

	go game.Q(proto1, K, 1,
		conn, "localhost", conn.Port(2*K), // P->
		conn, "localhost", conn.Port(0), // ->B
		conn, "localhost", conn.Port(K), // ->C
		wg)
	wg.Wait()
}
