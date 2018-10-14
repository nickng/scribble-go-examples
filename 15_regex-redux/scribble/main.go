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
 */

//go:generate scribblec-param.sh ../Regex.scr -d ../ -param Proto github.com/nickng/scribble-go-examples/15_regex-redux/Regex -param-api A -param-api B -param-api C

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/15_regex-redux/Regex/Proto"
	"github.com/nickng/scribble-go-examples/15_regex-redux/Regex/Proto/A_1to1"
	"github.com/nickng/scribble-go-examples/15_regex-redux/Regex/Proto/B_1toK"
	"github.com/nickng/scribble-go-examples/15_regex-redux/Regex/Proto/C_1to1"
	session "github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
)

var allvariants = []string{
	"agggtaaa|tttaccct",
	"[cgt]gggtaaa|tttaccc[acg]",
	"a[act]ggtaaa|tttacc[agt]t",
	"ag[act]gtaaa|tttac[agt]ct",
	"agg[act]taaa|ttta[agt]cct",
	"aggg[acg]aaa|ttt[cgt]ccct",
	"agggt[cgt]aa|tt[acg]accct",
	"agggta[cgt]a|t[acg]taccct",
	"agggtaa[cgt]|[acg]ttaccct",
	"agggtaaa|tttaccct",
	"agggtaaa|tttaccct",
	"[cgt]gggtaaa|tttaccc[acg]",
	"a[act]ggtaaa|tttacc[agt]t",
	"ag[act]gtaaa|tttac[agt]ct",
	"agg[act]taaa|ttta[agt]cct",
	"aggg[acg]aaa|ttt[cgt]ccct",
	"agggt[cgt]aa|tt[acg]accct",
	"agggta[cgt]a|t[acg]taccct",
	"agggtaa[cgt]|[acg]ttaccct",
	"[cgt]gggtaaa|tttaccc[acg]",
	"a[act]ggtaaa|tttacc[agt]t",
	"ag[act]gtaaa|tttac[agt]ct",
	"agg[act]taaa|ttta[agt]cct",
	"aggg[acg]aaa|ttt[cgt]ccct",
	"agggt[cgt]aa|tt[acg]accct",
	"agggta[cgt]a|t[acg]taccct",
	"agggtaa[cgt]|[acg]ttaccct",
}

type Subst struct {
	pat, repl string
}

var substs = []Subst{
	Subst{"B", "(c|g|t)"},
	Subst{"D", "(a|g|t)"},
	Subst{"H", "(a|c|t)"},
	Subst{"K", "(g|t)"},
	Subst{"M", "(a|c)"},
	Subst{"N", "(a|c|g|t)"},
	Subst{"R", "(a|g)"},
	Subst{"S", "(c|g)"},
	Subst{"V", "(a|c|g)"},
	Subst{"W", "(a|t)"},
	Subst{"Y", "(c|t)"},
}

func countMatches(pat string, bytes []byte) int {
	re := regexp.MustCompile(pat)
	n := 0
	for {
		e := re.FindIndex(bytes)
		if e == nil {
			break
		}
		n++
		bytes = bytes[e[1]:]
	}
	return n
}

var nCPU int

func main() {
	run_startt := time.Now()
	flag.IntVar(&nCPU, "ncpu", 8, "num goroutines")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	variants := allvariants[:nCPU]
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't read input: %s\n", err)
		os.Exit(2)
	}
	ilen := len(bytes)
	// Delete the comment lines and newlines
	bytes = regexp.MustCompile("(>[^\n]+)?\n").ReplaceAll(bytes, []byte{})
	clen := len(bytes)

	// Connections are created through shm.Listen(), instead of make(chan int).
	// They are generated once, and stored in a slice. The original code creates
	// the necessary channels just before running the corresponding worker.

	connB := make([]*shm.ShmListener, nCPU) // [1: B1-A, 2: B2-A, ...]
	for i := 0; i < nCPU; i++ {
		connB[i], err = shm.Listen(i + 1)
		defer connB[i].Close()
	}
	connC, err := shm.Listen(nCPU + 2) // [nCPU+2: C-A]
	defer connC.Close()

	prot := Proto.New()
	mini := prot.New_A_1to1(nCPU, 1)
	wg := new(sync.WaitGroup)
	wg.Add(len(connB) + 1)
	for i := range connB {
		go func(i int) {
			if err := mini.B_1toK_Accept(i+1, connB[i], new(session.PassByPointer)); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(i)
	}
	go func() {
		if err := mini.C_1to1_Accept(1, connC, new(session.PassByPointer)); err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	// main session initiated, main function created
	mmain := func() {
		mini.Run(master(ilen, clen, variants))
	}

	// initialise C session
	bb := bytes

	cini := prot.New_C_1to1(1)
	if err := cini.A_1to1_Dial(1, "_", nCPU+2, shm.Dial, new(session.PassByPointer)); err != nil {
		log.Fatal(err)
	}

	// C main function
	cmain := func() {
		cini.Run(substr(bb))
	}

	mkbmain := func(idx int) func() {
		bini := prot.New_B_1toK(nCPU, idx+1)
		if err := bini.A_1to1_Dial(1, "_", idx+1, shm.Dial, new(session.PassByPointer)); err != nil {
			log.Fatal(err)
		}
		return func() {
			bini.Run(worker(bytes))
		}
	}

	bmains := make([]func(), nCPU)
	for idx := 0; idx < nCPU; idx++ {
		bmains[idx] = mkbmain(idx)
	}
	wg.Wait()

	// Launch workers. Unlike in the first program, they stop at the first
	// receive until master distributes tasks.
	// In the original program, they start computing earlier, right after
	// channel creation.

	go cmain()
	for idx := 0; idx < nCPU; idx++ {
		go bmains[idx]()
	}

	mmain()
	run_endt := time.Now()
	fmt.Println(ilen, "\t", nCPU, "\t", run_endt.Sub(run_startt).Nanoseconds())
}

func substr(bb []byte) func(*C_1to1.Init) C_1to1.End {
	return func(s0 *C_1to1.Init) C_1to1.End {
		measures := make([]int, 1, 1)
		s1 := s0.A_1_Gather_Measure(measures)
		/*** Exactly as base program, after Measure() for synchronisation ***/
		for _, sub := range substs {
			bb = regexp.MustCompile(sub.pat).ReplaceAll(bb, []byte(sub.repl))
		}
		end := s1.A_1_Scatter_Len([]int{len(bb)})
		return *end
	}
}

func worker(bytes []byte) func(*B_1toK.Init) B_1toK.End {
	return func(s0 *B_1toK.Init) B_1toK.End {
		// Count receives variant, and calls countMatches, just as the original
		// program. The result is sent using Donec, instead of a custom channel.
		counts := make([]string, 1, 1)
		s1 := s0.A_1_Gather_Count(counts)
		end := s1.A_1_Scatter_Donec([]int{countMatches(counts[0], bytes)})
		return *end
	}
}

func master(ilen, clen int, variants []string) func(*A_1to1.Init) A_1to1.End {
	return func(s0 *A_1to1.Init) A_1to1.End {
		// Send variants through a channel. In base case, variants are passed
		// to goroutines as function arguments. The base case should be faster.
		s1 := s0.B_1toK_Scatter_Count(variants)

		// After workers received the interest variants,
		// measure sends a token to worker C to continue
		s2 := s1.C_1_Scatter_Measure([]int{0})

		// Wait for workers to finish and gather results.
		// Original program does not need this, since the receives are done
		// while printing results
		rs := make([]int, nCPU, nCPU)
		s3 := s2.B_1toK_Gather_Donec(rs)
		lens := make([]int, 1, 1)
		end := s3.C_1_Gather_Len(lens)

		for i, c := range rs {
			fmt.Fprintf(ioutil.Discard, "%s %d\n", variants[i], c)
		}
		fmt.Fprintf(ioutil.Discard, "\n%d\n%d\n%d\n", ilen, clen, lens[0])
		return *end
	}
}
