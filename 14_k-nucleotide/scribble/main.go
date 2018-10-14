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

//go:generate scribblec-param.sh ../KNuc.scr -d ../ -param Proto github.com/nickng/scribble-go-examples/14_k-nucleotide/KNuc -param-api A -param-api B -param-api S

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/14_k-nucleotide/KNuc/Proto"
	"github.com/nickng/scribble-go-examples/14_k-nucleotide/KNuc/Proto/A_1to1"
	"github.com/nickng/scribble-go-examples/14_k-nucleotide/KNuc/Proto/B_1toK"
	"github.com/nickng/scribble-go-examples/14_k-nucleotide/KNuc/Proto/S_1to2"
	session "github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
)

func count(data string, n int) map[string]int {
	counts := make(map[string]int)
	top := len(data) - n
	for i := 0; i <= top; i++ {
		s := data[i : i+n]
		counts[s]++
	}
	return counts
}

func countOne(data string, s string) int {
	return count(data, len(s))[s]
}

type kNuc struct {
	name  string
	count int
}

type kNucArray []kNuc

func (kn kNucArray) Len() int      { return len(kn) }
func (kn kNucArray) Swap(i, j int) { kn[i], kn[j] = kn[j], kn[i] }
func (kn kNucArray) Less(i, j int) bool {
	if kn[i].count == kn[j].count {
		return kn[i].name > kn[j].name // sort down
	}
	return kn[i].count > kn[j].count
}

func sortedArray(m map[string]int) kNucArray {
	kn := make(kNucArray, len(m))
	i := 0
	for k, v := range m {
		kn[i] = kNuc{k, v}
		i++
	}
	sort.Sort(kn)
	return kn
}

func printKnucs(a kNucArray) {
	sum := 0
	for _, kn := range a {
		sum += kn.count
	}
	for _, kn := range a {
		fmt.Fprintf(ioutil.Discard, "%s %.3f\n", kn.name, 100*float64(kn.count)/float64(sum))
	}
	fmt.Fprintf(ioutil.Discard, "\n")
}

var nCPU int

func main() {
	run_startt := time.Now()
	flag.IntVar(&nCPU, "ncpu", 8, "GOMAXPROCS")
	flag.Parse()
	runtime.GOMAXPROCS(8)
	in := bufio.NewReader(os.Stdin)
	three := []byte(">THREE ")
	for {
		line, err := in.ReadSlice('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "ReadLine err:", err)
			os.Exit(2)
		}
		if line[0] == '>' && bytes.Equal(line[0:len(three)], three) {
			break
		}
	}
	data, err := ioutil.ReadAll(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ReadAll err:", err)
		os.Exit(2)
	}
	// delete the newlines and convert to upper case
	j := 0
	for i := 0; i < len(data); i++ {
		if data[i] != '\n' {
			data[j] = data[i] &^ ' ' // upper case
			j++
		}
	}
	str := string(data[0:j])

	interests := []string{"GGT", "GGTA", "GGTATT", "GGTATTTTAATT",
		"GGT", "GGTA", "GGTATT", "GGTATTTTAATT",
		"GGT", "GGTA", "GGTATT", "GGTATTTTAATT",
		"GGT", "GGTA", "GGTATT", "GGTATTTTAATT",
		"GGT", "GGTA", "GGTATT", "GGTATTTTAATT"}

	var arr1, arr2 kNucArray

	// Create connections
	connS := make([]*shm.ShmListener, 2) // [1: S1-A, 2: S2-A]
	for i := 0; i < 2; i++ {
		connS[i], err = shm.Listen(1 + i)
		defer connS[i].Close()
	}
	connB := make([]*shm.ShmListener, nCPU) // [3: B1-A, 4: B2-A, ...]
	for i := 0; i < nCPU; i++ {
		connB[i], err = shm.Listen(3 + i)
		defer connB[i].Close()
	}

	// instantiate protocol
	prot := Proto.New()
	mini := prot.New_A_1to1(nCPU, 1)
	wg := new(sync.WaitGroup)
	wg.Add(len(connS) + len(connB))
	for i := range connS {
		go func(i int) {
			if err := mini.S_1to2_Accept(i+1, connS[i], new(session.PassByPointer)); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(i)
	}
	for i := range connB {
		go func(i int) {
			if err := mini.B_1toK_Accept(i+1, connB[i], new(session.PassByPointer)); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(i)
	}

	// main session initiated, main function created
	mmain := func() {
		mini.Run(master(&arr1, &arr2, interests))
	}

	inp := [2]*kNucArray{&arr1, &arr2}

	sorterInitialise := func(idx int) func() {
		ini := prot.New_S_1to2(idx + 1)
		ini.A_1to1_Dial(1, "_", 1+idx, shm.Dial, new(session.PassByPointer)) // A_idx-S
		return func() {
			ini.Run(sorter(idx, str, inp[idx]))
		}
	}
	sorter1 := sorterInitialise(0)
	sorter2 := sorterInitialise(1)

	workerInitialise := func(idx int) func() {
		ini := prot.New_B_1toK(nCPU, idx+1)
		ini.A_1to1_Dial(1, "_", 3+idx, shm.Dial, new(session.PassByPointer)) // B_idx-A
		return func() {
			ini.Run(worker(str))
		}
	}
	workers := make([]func(), nCPU)
	for idx := 0; idx < nCPU; idx++ {
		workers[idx] = workerInitialise(idx)
	}
	wg.Wait()

	// run sorters + workers
	go sorter1()
	go sorter2()
	for idx := 0; idx < nCPU; idx++ {
		go workers[idx]()
	}
	mmain()
	run_endt := time.Now()

	fmt.Println(len(data), "\t", nCPU, "\t", run_endt.Sub(run_startt).Nanoseconds())
}

func worker(str string) func(*B_1toK.Init) B_1toK.End {
	return func(s0 *B_1toK.Init) B_1toK.End {
		matches := make([]string, 1, 1)
		s1 := s0.A_1_Gather_Match(matches)
		result := []string{fmt.Sprintf("%d %s\n", countOne(str, matches[0]), matches[0])}
		end := s1.A_1_Scatter_Gather(result)
		return *end
	}
}

func sorter(i int, str string, arr *kNucArray) func(*S_1to2.Init) S_1to2.End {
	return func(s0 *S_1to2.Init) S_1to2.End {
		sorts := make([]int, 1, 1)
		s1 := s0.A_1_Gather_Sort(sorts)
		*arr = sortedArray(count(str, i+1))
		end := s1.A_1_Scatter_Done([]int{i})
		return *end
	}
}

func master(arr1, arr2 *kNucArray, interests []string) func(*A_1to1.Init) A_1to1.End {
	return func(s0 *A_1to1.Init) A_1to1.End {
		dones := make([]int, 2, 2)
		s2 := s0.
			S_1to2_Scatter_Sort([]int{1, 2}).
			B_1toK_Scatter_Match(interests[:nCPU]).
			S_1to2_Gather_Done(dones)
		printKnucs(*arr1)
		printKnucs(*arr2)
		rcs := make([]string, nCPU)
		end := s2.B_1toK_Gather_Gather(rcs)
		for _, rc := range rcs {
			fmt.Fprint(ioutil.Discard, rc)
		}
		return *end
	}
}
