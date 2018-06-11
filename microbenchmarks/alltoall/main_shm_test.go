package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Alltoall"
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Alltoall/A_1toM"
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Alltoall/B_1toN"
	"github.com/nickng/scribble-go-examples/microbenchmarks/message"
	session "github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
)

func initSM(t testing.TB, Count, M, N int) (roles []rolesAlltoall, cleanupFn func()) {
	prot := Alltoall.New()
	roles = make([]rolesAlltoall, Count)
	for c := 0; c < Count; c++ {
		roles[c].prot = prot
		roles[c].AM = make([]*A_1toM.A_1toM, M)
		roles[c].BN = make([]*B_1toN.B_1toN, N)
		// ---- Protocol initialisation ----
		for m := 0; m < M; m++ {
			roles[c].AM[m] = prot.New_A_1toM(N, M, m+1)
			roles[c].AM[m].Params["M"] = M
			roles[c].AM[m].Params["N"] = N
		}
		for n := 0; n < N; n++ {
			roles[c].BN[n] = prot.New_B_1toN(M, N, n+1)
			roles[c].BN[n].Params["M"] = M
			roles[c].BN[n].Params["N"] = N
		}
		// ---- Protocol initialisation END ----
	}
	var err error
	lnMN := make([][]*shm.Listener, M)
	for m := 0; m < M; m++ {
		lnMN[m] = make([]*shm.Listener, N)
		for n := 0; n < N; n++ {
			lnMN[m][n], err = shm.Listen(m*N + n)
			if err != nil {
				t.Error(err)
			}
		}
	}
	wg := new(sync.WaitGroup)
	wg.Add(M + N)
	for m := 0; m < M; m++ {
		go func(m int) {
			for c := 0; c < Count; c++ {
				if c == 0 {
					for n := 0; n < N; n++ {
						if err := roles[c].AM[m].B_1toN_Dial(n+1, "", m*N+n, shm.Dial, new(session.PassByPointer)); err != nil {
							t.Error(err)
						}
						roles[c].AM[m].MPChan.CheckConnection()
					}
				} else {
					for role, conns := range roles[c-1].AM[m].MPChan.Conns {
						for id, conn := range conns {
							roles[c].AM[m].MPChan.Conns[role][id] = conn
						}
					}
					for role, fmts := range roles[c-1].AM[m].MPChan.Fmts {
						for id, f := range fmts {
							roles[c].AM[m].MPChan.Fmts[role][id] = f
						}
					}
				}
			}
			wg.Done()
		}(m)
	}
	for n := 0; n < N; n++ {
		go func(n int) {
			for c := 0; c < Count; c++ {
				if c == 0 {
					for m := 0; m < M; m++ {
						if err := roles[c].BN[n].A_1toM_Accept(m+1, lnMN[m][n], new(session.PassByPointer)); err != nil {
							t.Error(err)
						}
					}
					roles[c].BN[n].MPChan.CheckConnection()
				} else {
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
		for _, lnN := range lnMN {
			for _, ln := range lnN {
				ln.Close()
			}
		}
	}
	return roles, cleanupFn
}

func manualSM(t testing.TB, Count, M, N int) {
	b, benchmarkMode := t.(*testing.B)
	roles, cleanupFn := initSM(t, Count, M, N)
	wg := new(sync.WaitGroup)
	wg.Add(M + N)
	wgStart := new(sync.WaitGroup)
	wgStart.Add(M + N)
	wgEnd := new(sync.WaitGroup)
	wgEnd.Add(M + N)
	for i := 0; i < M; i++ {
		go func(m int) {
			vals := make([]message.Int, N)
			wgStart.Done()
			wgStart.Wait()
			if benchmarkMode {
				b.ResetTimer()
			}
			for c := 0; c < Count; c++ {
				// ---- Begin overhead measurement ----
				for n := 0; n < N; n++ {
					roles[c].AM[m].MPChan.Fmts["B"][n+1].Serialize(&vals[n])
				}
				// ---- End overhead measurement ----
			}
			wgEnd.Done()
			wgEnd.Wait()
			wg.Done()
		}(i)
	}
	for j := 0; j < N; j++ {
		go func(n int) {
			val := make([]message.Int, M)
			wgStart.Done()
			wgStart.Wait()
			for c := 0; c < Count; c++ {
				// ---- Begin overhead measurement ----
				for m := 0; m < M; m++ {
					var v session.ScribMessage = &val[m]
					roles[c].BN[n].MPChan.Fmts["A"][m+1].Deserialize(&v)
				}
				// ---- End overhead measurement ----
			}
			wgEnd.Done()
			wgEnd.Wait()
			wg.Done()
		}(j)
	}
	wg.Wait()
	cleanupFn()
}

func scribbleSM(t testing.TB, Count, M, N int) {
	b, benchmarkMode := t.(*testing.B)
	roles, cleanupFn := initSM(t, Count, M, N)
	wg := new(sync.WaitGroup)
	wg.Add(M + N)
	wgStart := new(sync.WaitGroup)
	wgStart.Add(M + N)
	wgEnd := new(sync.WaitGroup)
	wgEnd.Add(M + N)
	for i := 0; i < M; i++ {
		go func(m int) {
			vals := make([]message.Int, N)
			wgStart.Done()
			wgStart.Wait()
			if m == 0 && benchmarkMode {
				b.ResetTimer()
			}
			for c := 0; c < Count; c++ {
				// ---- Begin overhead measurement ----
				roles[c].AM[m].Run(func(s *A_1toM.Init) A_1toM.End {
					return *(s.B_1toN_Scatter_Int(vals))
				})
				// ---- End overhead measurement ----
			}
			wgEnd.Done()
			wgEnd.Wait()
			if m == 0 && benchmarkMode {
				b.StopTimer()
			}
			wg.Done()
		}(i)
	}
	for j := 0; j < N; j++ {
		go func(n int) {
			xs := make([]message.Int, M)
			wgStart.Done()
			wgStart.Wait()
			for c := 0; c < Count; c++ {
				// ---- Begin overhead measurement ----
				roles[c].BN[n].Run(func(s *B_1toN.Init) B_1toN.End {
					return *(s.A_1toM_Gather_Int(xs))
				})
				// ---- End overhead measurement ----
			}
			wgEnd.Done()
			wgEnd.Wait()
			wg.Done()
		}(j)
	}
	wg.Wait()
	cleanupFn()
}

func TestManualSM(t *testing.T) {
	for m := MinM; m <= MaxM; m++ {
		t.Run(fmt.Sprintf("M=%d", m), func(t *testing.T) {
			for n := MinN; n <= MaxN; n++ {
				t.Run(fmt.Sprintf("N=%d", n), func(t *testing.T) {
					manualSM(t, 1, m, n)
				})
			}
		})
	}
}

func TestScribbleSM(t *testing.T) {
	for m := MinM; m <= MaxM; m++ {
		t.Run(fmt.Sprintf("M=%d", m), func(t *testing.T) {
			for n := MinN; n <= MaxN; n++ {
				t.Run(fmt.Sprintf("N=%d", n), func(t *testing.T) {
					scribbleSM(t, 1, m, n)
				})
			}
		})
	}
}

func BenchmarkManualSM(b *testing.B) {
	k := 0
	for m, n := MinM, MinN; m <= MaxM && n <= MaxN; m, n = m+(k%2), n+((k+1)%2) {
		k++
		b.Run(fmt.Sprintf("N=%d", m+n), func(b *testing.B) {
			manualSM(b, b.N, m, n)
		})
	}
}

func BenchmarkScribbleSM(b *testing.B) {
	k := 0
	for m, n := MinM, MinN; m <= MaxM && n <= MaxN; m, n = m+(k%2), n+((k+1)%2) {
		k++
		b.Run(fmt.Sprintf("N=%d", m+n), func(b *testing.B) {
			scribbleSM(b, b.N, m, n)
		})
	}
}
