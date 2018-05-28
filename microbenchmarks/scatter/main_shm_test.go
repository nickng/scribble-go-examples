package main_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Scatter"
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Scatter/A_1to1"
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Scatter/B_1toN"
	"github.com/nickng/scribble-go-examples/microbenchmarks/message"
	session "github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
)

func initSM(t testing.TB, Count, N int) (roles []rolesScatter, cleanupFn func()) {
	prot := Scatter.New()
	roles = make([]rolesScatter, Count)
	for c := 0; c < Count; c++ {
		// ---- Protocol initialisation ----
		roles[c].prot = prot
		roles[c].A = prot.New_A_1to1(N, 1)
		roles[c].A.Params["N"] = N
		roles[c].BN = make([]*B_1toN.B_1toN, N)
		for n := 0; n < N; n++ {
			roles[c].BN[n] = prot.New_B_1toN(N, n+1)
			roles[c].BN[n].Params["N"] = N
		}
		// ---- Protocol initialisation END ----
	}
	var err error
	lnN := make([][]*shm.ShmListener, Count)
	for c := 0; c < Count; c++ {
		lnN[c] = make([]*shm.ShmListener, N)
		for n := 0; n < N; n++ {
			lnN[c][n], err = shm.Listen(c*N + n)
			if err != nil {
				t.Error(err)
			}
		}
	}
	wg := new(sync.WaitGroup)
	wg.Add(1 + N)
	go func() {
		for c := 0; c < Count; c++ {
			if c == 0 {
				for n := 0; n < N; n++ {
					if err := roles[c].A.B_1toN_Dial(1+n, "", c*N+n, shm.Dial, new(session.PassByPointer)); err != nil {
						t.Error(err)
					}
				}
				roles[c].A.MPChan.CheckConnection()
			} else {
				// Reuse connection
				for role, conns := range roles[c-1].A.MPChan.Conns {
					for id, conn := range conns {
						roles[c].A.MPChan.Conns[role][id] = conn
					}
				}
				for role, fmts := range roles[c-1].A.MPChan.Fmts {
					for id, f := range fmts {
						roles[c].A.MPChan.Fmts[role][id] = f
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
					if err := roles[c].BN[n].A_1to1_Accept(1, lnN[c][n], new(session.PassByPointer)); err != nil {
						t.Error(err)
					}
					roles[c].BN[n].MPChan.CheckConnection()
				} else {
					// Reuse connection
					for role, conns := range roles[c-1].BN[n].MPChan.Conns {
						for id, conn := range conns {
							roles[c].BN[n].MPChan.Conns[role][id] = conn
						}
					}
					for role, fmts := range roles[c-1].BN[n].MPChan.Fmts {
						for id, f := range fmts {
							roles[c].BN[n].MPChan.Fmts[role][id] = f
						}
					}
				}
			}
			wg.Done()
		}(n)
	}
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
	wg.Add(N + 1)
	wgStart := new(sync.WaitGroup)
	wgStart.Add(N + 1)
	wgEnd := new(sync.WaitGroup)
	wgEnd.Add(N + 1)
	go func() {
		var v message.Int
		wgStart.Done()
		wgStart.Wait()
		if benchmarkMode {
			b.ResetTimer()
		}
		for c := 0; c < Count; c++ {
			// ---- Begin overhead measurement ----
			for n := 0; n < N; n++ {
				roles[c].A.MPChan.Fmts["B"][n+1].Serialize(&v)
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
	for i := 0; i < N; i++ {
		go func(n int) {
			var v session.ScribMessage = new(message.Int)
			wgStart.Done()
			wgStart.Wait()
			for c := 0; c < Count; c++ {
				// ---- Begin overhead measurement ----
				roles[c].BN[n].MPChan.Fmts["A"][1].Deserialize(&v)
				// ---- End overhead measurement ----
			}
			wgEnd.Done()
			wgEnd.Wait()
			wg.Done()
		}(i)
	}
	wg.Wait()
	cleanupFn()
}

func scribbleSM(t testing.TB, Count, N int) {
	b, benchmarkMode := t.(*testing.B)
	roles, cleanupFn := initSM(t, Count, N)
	wg := new(sync.WaitGroup)
	wg.Add(N + 1)
	wgStart := new(sync.WaitGroup)
	wgStart.Add(N + 1)
	wgEnd := new(sync.WaitGroup)
	wgEnd.Add(N + 1)
	go func() {
		vals := make([]message.Int, N)
		wgStart.Done()
		wgStart.Wait()
		if benchmarkMode {
			b.ResetTimer()
		}
		for c := 0; c < Count; c++ {
			// ---- Begin overhead measurement ----
			roles[c].A.Run(func(s *A_1to1.Init) A_1to1.End {
				return *(s.B_1toN_Scatter_Int(vals))
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
	for i := 0; i < N; i++ {
		go func(n int) {
			vals := make([]message.Int, 1)
			wgStart.Done()
			wgStart.Wait()
			for c := 0; c < Count; c++ {
				// ---- Begin overhead measurement ----
				roles[c].BN[n].Run(func(s *B_1toN.Init) B_1toN.End {
					return *(s.A_1to1_Gather_Int(vals))
				})
				// ---- End overhead measurement ----
			}
			wgEnd.Done()
			wgEnd.Wait()
			wg.Done()
		}(i)
	}
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
