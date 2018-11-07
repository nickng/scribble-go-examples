package main_test

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Gather"
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Gather/A_1toN"
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Gather/B_1to1"
	"github.com/nickng/scribble-go-examples/microbenchmarks/message"
	session "github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"
)

func init() {
	gob.Register(&message.Int{})
}

func initTCP(t testing.TB, Count, N int) (roles []rolesGather, cleanupFn func()) {
	prot := Gather.New()
	roles = make([]rolesGather, Count)
	for c := 0; c < Count; c++ {
		// ---- Protocol initialisation ----
		roles[c].prot = prot
		roles[c].AN = make([]*A_1toN.A_1toN, N)
		for n := 0; n < N; n++ {
			roles[c].AN[n] = prot.New_A_1toN(N, n+1)
			roles[c].AN[n].Params["N"] = N
		}
		roles[c].B = prot.New_B_1to1(N, 1)
		roles[c].B.Params["N"] = N
		// ---- Protocol initialisation END ----
	}
	var err error
	lnN := make([]*tcp.TcpListener, N)
	for n := 0; n < N; n++ {
		lnN[n], err = tcp.Listen(0) // Auto select port.
		if err != nil {
			t.Error(err)
		}
	}
	wg := new(sync.WaitGroup)
	wg.Add(1 + N)
	go func() {
		for c := 0; c < Count; c++ {
			if c == 0 {
				for n := 0; n < N; n++ {
					addr := lnN[n].Addr().(*net.TCPAddr)
					if err := roles[c].B.A_1toN_Dial(1+n, addr.IP.String(), addr.Port, tcp.Dial, new(session.GobFormatter)); err != nil {
						panic(err)
					}
				}
				roles[c].B.MPChan.CheckConnection()
			} else {
				// Reuse connection
				for role, conns := range roles[c-1].B.MPChan.Conns {
					for id, conn := range conns {
						roles[c].B.MPChan.Conns[role][id] = conn
					}
				}
				for role, fmts := range roles[c-1].B.MPChan.Fmts {
					for id, f := range fmts {
						roles[c].B.MPChan.Fmts[role][id] = f
					}
				}
			}
		}
		wg.Done()
	}()
	for n := 0; n < N; n++ {
		go func(n int) {
			for c := 0; c < Count; c++ {
				if c == 0 {
					if err := roles[c].AN[n].B_1to1_Accept(1, lnN[n], new(session.GobFormatter)); err != nil {
						panic(err)
					}
					roles[c].AN[n].MPChan.CheckConnection()
				} else {
					// Reuse connection
					for role, conns := range roles[c-1].AN[n].MPChan.Conns {
						for id, conn := range conns {
							roles[c].AN[n].MPChan.Conns[role][id] = conn
						}
					}
					for role, fmts := range roles[c-1].AN[n].MPChan.Fmts {
						for id, f := range fmts {
							roles[c].AN[n].MPChan.Fmts[role][id] = f
						}
					}
				}
			}
			wg.Done()
		}(n)
	}
	wg.Wait()
	cleanupFn = func() {
		roles[0].B.MPChan.Close()
		for n := 0; n < N; n++ {
			roles[0].AN[n].MPChan.Close()
		}
		for _, ln := range lnN {
			if err := ln.Close(); err != nil {
				t.Error(err)
			}
		}
	}
	return roles, cleanupFn
}

func manualTCP(t testing.TB, Count, N int) {
	b, benchmarkMode := t.(*testing.B)
	roles, cleanupFn := initTCP(t, Count, N)
	wg := new(sync.WaitGroup)
	wg.Add(N + 1)
	wgStart := new(sync.WaitGroup)
	wgStart.Add(N + 1)
	wgEnd := new(sync.WaitGroup)
	wgEnd.Add(N + 1)
	for i := 0; i < N; i++ {
		go func(n int) {
			var v message.Int
			wgStart.Done()
			wgStart.Wait()
			for c := 0; c < Count; c++ {
				// ---- Begin overhead measurement ----
				roles[c].AN[n].MPChan.Fmts["B"][1].Serialize(&v)
				// ---- End overhead measurement ----
			}
			wgEnd.Done()
			wgEnd.Wait()
			wg.Done()
		}(i)
	}
	go func() {
		var v session.ScribMessage = new(message.Int)
		wgStart.Done()
		wgStart.Wait()
		if benchmarkMode {
			b.ResetTimer()
		}
		for c := 0; c < Count; c++ {
			// ---- Begin overhead measurement ----
			for n := 0; n < N; n++ {
				roles[c].B.MPChan.Fmts["A"][n+1].Deserialize(&v)
			}
			// ---- End overhead measurement ----
		}
		wgEnd.Done()
		wgEnd.Wait()
		if benchmarkMode {
			b.StopTimer()
		}
		wg.Done()
	}()
	wg.Wait()
	cleanupFn()
}

func scribbleTCP(t testing.TB, Count, N int) {
	b, benchmarkMode := t.(*testing.B)
	roles, cleanupFn := initTCP(t, Count, N)
	wg := new(sync.WaitGroup)
	wg.Add(1 + N)
	wgStart := new(sync.WaitGroup)
	wgStart.Add(1 + N)
	wgEnd := new(sync.WaitGroup)
	wgEnd.Add(1 + N)
	for i := 0; i < N; i++ {
		go func(n int) {
			vals := make([]message.Int, 1)
			wgStart.Done()
			wgStart.Wait()
			for c := 0; c < Count; c++ {
				// ---- Begin overhead measurement ----
				roles[c].AN[n].Run(func(s *A_1toN.Init) A_1toN.End {
					return *(s.B_1_Scatter_Int(vals))
				})
				// ---- End overhead measurement ----
			}
			wgEnd.Done()
			wgEnd.Wait()
			wg.Done()
		}(i)
	}
	go func() {
		vals := make([]message.Int, N)
		wgStart.Done()
		wgStart.Wait()
		if benchmarkMode {
			b.ResetTimer()
		}
		for c := 0; c < Count; c++ {
			// ---- Begin overhead measurement ----
			roles[c].B.Run(func(s *B_1to1.Init) B_1to1.End {
				return *(s.A_1toN_Gather_Int(vals))
			})
			// ---- End overhead measurement ----
		}
		wgEnd.Done()
		wgEnd.Wait()
		if benchmarkMode {
			b.StopTimer()
		}
		wg.Done()
	}()
	wg.Wait()
	cleanupFn()
}

func TestNativeTCP(t *testing.T) {
	for i := MinN; i <= MaxN; i++ {
		t.Run(fmt.Sprintf("N=%d", i), func(t *testing.T) {
			manualTCP(t, 1, i)
		})
	}
}

func TestScribbleTCP(t *testing.T) {
	for i := MinN; i <= MaxN; i++ {
		t.Run(fmt.Sprintf("N=%d", i), func(t *testing.T) {
			scribbleTCP(t, 1, i)
		})
	}
}

func BenchmarkNativeTCP(b *testing.B) {
	for i := MinN; i <= MaxN; i++ {
		b.Run(fmt.Sprintf("N=%d", i), func(b *testing.B) {
			manualTCP(b, b.N, i)
		})
	}
}

func BenchmarkScribbleTCP(b *testing.B) {
	for i := MinN; i <= MaxN; i++ {
		b.Run(fmt.Sprintf("N=%d", i), func(b *testing.B) {
			scribbleTCP(b, b.N, i)
		})
	}
}
