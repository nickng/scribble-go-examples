/*
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

    * Redistributions of source code must retain the above copyright
    notice, this list of conditions and the following disclaimer.

    * Redistributions in binary form must reproduce the above copyright
    notice, this list of conditions and the following disclaimer in the
    documentation and/or other materials provided with the distribution.

    * Neither the name of "The Computer Language Benchmarks Game" nor the
    name of "The Computer Language Shootout Benchmarks" nor the names of
    its contributors may be used to endorse or promote products derived
    from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED.  IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/

/* The Computer Language Benchmarks Game
 * http://shootout.alioth.debian.org/
 *
 * contributed by The Go Authors.
 * Based on spectral-norm.c by Sebastien Loisel
 */

//go:generate scribblec-param.sh ../SN.scr -d ../ -param Proto github.com/nickng/scribble-go-examples/16_spectral-norm/SN -param-api A -param-api B

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"

	"runtime"
	"time"

	"github.com/nickng/scribble-go-examples/16_spectral-norm/SN/Proto"
	"github.com/nickng/scribble-go-examples/16_spectral-norm/SN/Proto/A_1to1"
	"github.com/nickng/scribble-go-examples/16_spectral-norm/SN/Proto/B_1toK"
	session "github.com/rhu1/scribble-go-runtime/runtime/session2"
	shm "github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
)

var n = flag.Int("n", 2000, "count")
var nCPU = flag.Int("ncpu", 4, "number of cpus")

func evalA(i, j int) float64 { return 1 / float64(((i+j)*(i+j+1)/2 + i + 1)) }

type Vec []float64

func (v Vec) Times(i, n int, u Vec) {
	for ; i < n; i++ {
		v[i] = 0
		for j := 0; j < len(u); j++ {
			v[i] += evalA(i, j) * u[j]
		}
	}
}

func (v Vec) TimesTransp(i, n int, u Vec) {
	for ; i < n; i++ {
		v[i] = 0
		for j := 0; j < len(u); j++ {
			v[i] += evalA(j, i) * u[j]
		}
	}
}

var err error

func main() {
	run_startt := time.Now()
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	N := *n
	u := make(Vec, N)
	for i := 0; i < N; i++ {
		u[i] = 1
	}
	v := make(Vec, N)

	conns := make([]*shm.ShmListener, *nCPU) // [1: A-B1, 2: A-B2, ...]
	for i := range conns {
		conns[i], err = shm.Listen(i + 1)
		defer conns[i].Close()
	}

	// instantiate protocol
	prot := Proto.New()
	mini := prot.New_A_1to1(*nCPU, 1)
	for i, conn := range conns {
		if err := mini.B_1toK_Accept(i+1, conn, new(session.PassByPointer)); err != nil {
			log.Fatal(err)
		}
	}

	var x Vec
	// main session initiated, main function created
	mmain := func() {
		mini.Run(master(N, u, v, &x))
	}

	// instantiate workers with sub-roles
	workerInitialise := func(idx int) func() {
		ini := prot.New_B_1toK(*nCPU, idx+1)
		if err := ini.A_1to1_Dial(1, "_", idx+1, shm.Dial, new(session.PassByPointer)); err != nil {
			log.Fatal(err)
		}
		return func() {
			ini.Run(worker(idx, u, v, &x))
		}
	}

	workers := make([]func(), *nCPU)
	for i := 0; i < *nCPU; i++ {
		workers[i] = workerInitialise(i)
	}
	for i := 0; i < *nCPU; i++ {
		go workers[i]()
	}

	mmain()
	run_endt := time.Now()
	fmt.Println(N, "\t", *nCPU, "\t", run_endt.Sub(run_startt).Nanoseconds())
}

func worker(i int, u, v Vec, x *Vec) func(*B_1toK.Init) B_1toK.End {
	return func(s0 *B_1toK.Init) B_1toK.End {
		var buf = []int{0} // container for received message
		var pl int
		for {
			switch s1 := s0.A_1_Branch().(type) {
			case *B_1toK.Times:
				// x.Times u followed by v.TimesTransp x
				s2 := s1.Recv_Times(&pl)
				(*x).Times(i*len(v) / *nCPU, (i+1)*len(v) / *nCPU, u)
				// tell master we are done
				s3 := s2.
					A_1_Scatter_Done(buf).
					A_1_Gather_Next(buf)
				v.TimesTransp(i*len(v) / *nCPU, (i+1)*len(v) / *nCPU, *x)
				s4 := s3.
					A_1_Scatter_Done(buf).
					A_1_Gather_TimeStr(buf)

				// now we are doing a u.TimesTransp(v),
				// so u and v should be reversed in the operations.
				// Also, x should be fresh and of length (v).

				(*x).Times(i*len(u) / *nCPU, (i+1)*len(u) / *nCPU, v)
				s5 := s4.
					A_1_Scatter_Done(buf).
					A_1_Gather_Next(buf)
				u.TimesTransp(i*len(u) / *nCPU, (i+1)*len(u) / *nCPU, *x)
				s0 = s5.
					A_1_Scatter_Done(buf)
			case *B_1toK.Finish:
				end := s1.Recv_Finish(&pl)
				// last iteration
				return *end
			}
		}
	}
}

func master(N int, u, v Vec, x *Vec) func(*A_1to1.Init) A_1to1.End {
	return func(s0 *A_1to1.Init) A_1to1.End {
		buf := make([]int, *nCPU)
		pl := make([]int, *nCPU)
		for i := 0; i < *nCPU; i++ {
			pl[i] = 1 + i
		}
		for i := 0; i < 10; i++ {
			// v.ATimesTransp(u)
			*x = make(Vec, len(u))
			s1 := s0.
				B_1toK_Scatter_Times(pl).
				B_1toK_Gather_Done(buf).
				B_1toK_Scatter_Next(pl)
			// u.ATimesTransp(v)
			*x = make(Vec, len(v))
			s0 = s1.
				B_1toK_Gather_Done(buf).
				B_1toK_Scatter_TimeStr(pl).
				B_1toK_Gather_Done(buf).
				B_1toK_Scatter_Next(pl).
				B_1toK_Gather_Done(buf)
		}

		// after synchronisation finishes, continue with local computation
		var vBv, vv float64
		for i := 0; i < N; i++ {
			vBv += u[i] * v[i]
			vv += v[i] * v[i]
		}
		fmt.Fprintf(ioutil.Discard, "%0.9f\n", math.Sqrt(vBv/vv))

		// finalise
		end := s0.B_1toK_Scatter_Finish(pl)
		return *end
	}
}
