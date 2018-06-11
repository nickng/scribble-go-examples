package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Alltoall"
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Alltoall/A_1toM"
	"github.com/nickng/scribble-go-examples/microbenchmarks/Microbenchmarks/Alltoall/B_1toN"
	"github.com/nickng/scribble-go-examples/microbenchmarks/message"
	session "github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"
)

func init() {
	gob.Register(&message.Int{})
}

func initTCP(t testing.TB, Count, M, N int) (roles []rolesAlltoall, cleanupFn func()) {
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
	lnMN := make([][]*tcp.TcpListener, M)
	for m := 0; m < M; m++ {
		lnMN[m] = make([]*tcp.TcpListener, N)
		for n := 0; n < N; n++ {
			lnMN[m][n], err = tcp.Listen(0) // Auto select port.
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
						addr := lnMN[m][n].Addr().(*net.TCPAddr)
						if err := roles[c].AM[m].B_1toN_Dial(1+n, addr.IP.String(), addr.Port, tcp.Dial, new(session.GobFormatter)); err != nil {
							panic(err)
						}
					}
					roles[c].AM[m].MPChan.CheckConnection()
				} else {
					// reuse connection
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
						if err := roles[c].BN[n].A_1toM_Accept(m+1, lnMN[m][n], new(session.GobFormatter)); err != nil {
							panic(err)
						}
					}
					roles[c].BN[n].MPChan.CheckConnection()
				} else {
					// reuse connection
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
		for m := 0; m < M; m++ {
			roles[0].AM[m].MPChan.Close()
		}
		for n := 0; n < N; n++ {
			roles[0].BN[n].MPChan.Close()
		}
		for _, lnN := range lnMN {
			for _, ln := range lnN {
				if err := ln.Close(); err != nil {
					t.Error(err)
				}
			}
		}
	}
	return roles, cleanupFn
}

func manualTCP(t testing.TB, Count, M, N int) {
	b, benchmarkMode := t.(*testing.B)
	roles, cleanupFn := initTCP(t, Count, M, N)
	wg := new(sync.WaitGroup)
	wg.Add(M + N)
	wgStart := new(sync.WaitGroup)
	wgStart.Add(M + N)
	wgEnd := new(sync.WaitGroup)
	wgEnd.Add(M + N)
	for i := 0; i < M; i++ {
		go func(m int) {
			var v message.Int
			wgStart.Done()
			wgStart.Wait()
			if m == 0 && benchmarkMode {
				b.ResetTimer()
			}
			for c := 0; c < Count; c++ {
				// ---- Begin overhead measurement ----
				for n := 0; n < N; n++ {
					roles[c].AM[m].MPChan.Fmts["B"][n+1].Serialize(&v)
				}
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
			vs := make([]message.Int, M)
			wgStart.Done()
			wgStart.Wait()
			for c := 0; c < Count; c++ {
				// ---- Begin overhead measurement ----
				for m := 0; m < M; m++ {
					var v session.ScribMessage = &vs[m]
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

func scribbleTCP(t testing.TB, Count, M, N int) {
	b, benchmarkMode := t.(*testing.B)
	roles, cleanupFn := initTCP(t, Count, M, N)
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
			vals := make([]message.Int, M)
			wgStart.Done()
			wgStart.Wait()
			for c := 0; c < Count; c++ {
				// ---- Begin overhead measurement ----
				roles[c].BN[n].Run(func(s *B_1toN.Init) B_1toN.End {
					return *(s.A_1toM_Gather_Int(vals))
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

func TestManualTCP(t *testing.T) {
	for m := MinM; m <= MaxM; m++ {
		t.Run(fmt.Sprintf("M=%d", m), func(t *testing.T) {
			for n := MinN; n <= MaxN; n++ {
				t.Run(fmt.Sprintf("N=%d", n), func(t *testing.T) {
					manualTCP(t, 1, m, n)
				})
			}
		})
	}
}

func TestScribbleTCP(t *testing.T) {
	for m := MinM; m <= MaxM; m++ {
		t.Run(fmt.Sprintf("M=%d", m), func(t *testing.T) {
			for n := MinN; n <= MaxN; n++ {
				t.Run(fmt.Sprintf("N=%d", n), func(t *testing.T) {
					scribbleTCP(t, 1, m, n)
				})
			}
		})
	}
}

func BenchmarkManualTCP(b *testing.B) {
	k := 0
	for m, n := MinM, MinN; m <= MaxM && n <= MaxN; m, n = m+(k%2), n+((k+1)%2) {
		k++
		b.Run(fmt.Sprintf("N=%d", m+n), func(b *testing.B) {
			manualTCP(b, b.N, m, n)
		})
	}
}

func BenchmarkScribbleTCP(b *testing.B) {
	k := 0
	for m, n := MinM, MinN; m <= MaxM && n <= MaxN; m, n = m+(k%2), n+((k+1)%2) {
		k++
		b.Run(fmt.Sprintf("N=%d", m+n), func(b *testing.B) {
			scribbleTCP(b, b.N, m, n)
		})
	}
}
