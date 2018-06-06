package main_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Gather"
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Gather/A_1toN"
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Gather/B_1to1"
	"github.com/nickng/scribble-go-examples/microbenchmarks/message"
	session "github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
)

func initSM(t testing.TB, Count, N int) (roles []rolesGather, cleanupFn func()) {
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
	lnN := make([][]*shm.Listener, Count)
	for c := 0; c < Count; c++ {
		lnN[c] = make([]*shm.Listener, N)
		for n := 0; n < N; n++ {
			lnN[c][n], err = shm.Listen(c*N + n)
			if err != nil {
				t.Error(err)
			}
		}
	}
	wg := new(sync.WaitGroup)
	wg.Add(N + 1)
	for n := 0; n < N; n++ {
		go func(n int) {
			for c := 0; c < Count; c++ {
				if err := roles[c].AN[n].B_1to1_Accept(1, lnN[c][n], new(session.PassByPointer)); err != nil {
					t.Error(err)
				}
				roles[c].AN[n].MPChan.CheckConnection()
			}
			wg.Done()
		}(n)
	}
	go func() {
		for c := 0; c < Count; c++ {
			for n := 0; n < N; n++ {
				if err := roles[c].B.A_1toN_Dial(1+n, "", c*N+n, shm.Dial, new(session.PassByPointer)); err != nil {
					t.Error(err)
				}
			}
			roles[c].B.MPChan.CheckConnection()
		}
		wg.Done()
	}()
	wg.Wait()
	cleanupFn = func() {
		for _, ln := range lnN {
			for _, l := range ln {
				l.Close()
			}
		}
	}
	return roles, cleanupFn
}

func manualSM(t testing.TB, Count, N int) {
	b, benchmarkMode := t.(*testing.B)
	roles, cleanupFn := initSM(t, Count, N)
	wg := new(sync.WaitGroup)
	wg.Add(1 + N)
	wgStart := new(sync.WaitGroup)
	wgStart.Add(1 + N)
	wgEnd := new(sync.WaitGroup)
	wgEnd.Add(1 + N)
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

func scribbleSM(t testing.TB, Count, N int) {
	b, benchmarkMode := t.(*testing.B)
	roles, cleanupFn := initSM(t, Count, N)
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
					return *(s.B_1to1_Scatter_Int(vals))
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

func TestManualSM(t *testing.T) {
	for i := MinN; i <= MaxN; i++ {
		t.Run(fmt.Sprintf("N=%d", i), func(t *testing.T) {
			manualSM(t, 1, i)
		})
	}
}

func TestScribbleSM(t *testing.T) {
	for i := MinN; i <= MaxN; i++ {
		t.Run(fmt.Sprintf("N=%d", i), func(t *testing.T) {
			scribbleSM(t, 1, i)
		})
	}
}

func BenchmarkManualSM(b *testing.B) {
	for i := MinN; i <= MaxN; i++ {
		b.Run(fmt.Sprintf("N=%d", i), func(b *testing.B) {
			manualSM(b, b.N, i)
		})
	}
}

func BenchmarkScribbleSM(b *testing.B) {
	for i := MinN; i <= MaxN; i++ {
		b.Run(fmt.Sprintf("N=%d", i), func(b *testing.B) {
			scribbleSM(b, b.N, i)
		})
	}
}
