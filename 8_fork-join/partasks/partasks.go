//go:generate scribblec-param.sh ../ForkJoin.scr -d ../ -param Protocol github.com/nickng/scribble-go-examples/8_fork-join/ForkJoin -param-api Master -param-api Worker

// Package partasks implements the fork-join pattern.
// The name partasks stands for parallel tasks and is used to avoid
// name clash with the forkjoin protocol name.
package partasks

import (
	"encoding/gob"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/8_fork-join/ForkJoin/Protocol"
	"github.com/nickng/scribble-go-examples/8_fork-join/ForkJoin/Protocol/Master_1to1"
	"github.com/nickng/scribble-go-examples/8_fork-join/ForkJoin/Protocol/Worker_1toK"
	"github.com/nickng/scribble-go-examples/8_fork-join/message"
	"github.com/nickng/scribble-go-examples/scributil"
)

func init() {
	gob.Register(new(message.Map))
	gob.Register(new(message.Red))
}

// Master is the coordinator.
func Master(p *Protocol.Protocol, K, self int, cc scributil.ClientConn, host string, baseport int, wg *sync.WaitGroup) {
	Master := p.New_Master_1to1(K, self)
	wgCli := new(sync.WaitGroup)
	wgCli.Add(K)
	mu := new(sync.Mutex) // for concurrent access to connection map
	for k := 1; k <= K; k++ {
		go func(k int) {
			mu.Lock()
			if err := Master.Worker_1toK_Dial(k, host, baseport+k, cc.Dial, cc.Formatter()); err != nil {
				log.Fatalf("cannot dial: %v", err)
			}
			mu.Unlock()
			wgCli.Done()
		}(k)
	}
	wgCli.Wait()
	Master.Run(func(s *Master_1to1.Init) Master_1to1.End {
		var m []message.Map
		for i := 1; i <= K; i++ {
			m = append(m, message.Map{V: i})
		}
		s0 := s.Worker_1toK_Scatter_Map(m)
		scributil.Debugf("Master: sent %v.\n", m)
		// collect r
		r := make([]message.Red, K)
		sEnd := s0.Worker_1toK_Gather_Red(r)
		scributil.Debugf("Master: received %v.\n", r)
		return *sEnd
	})
	wg.Done()
}

// Worker is the task worker.
func Worker(p *Protocol.Protocol, K, self int, sc scributil.ServerConn, port int, wg *sync.WaitGroup) {
	Worker := p.New_Worker_1toK(K, self)
	ln, err := sc.Listen(port)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	if err := Worker.Master_1to1_Accept(1, ln, sc.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	Worker.Run(func(s *Worker_1toK.Init) Worker_1toK.End {
		m := make([]message.Map, 1)
		s0 := s.Master_1_Gather_Map(m)
		scributil.Debugf("Worker[%d]: received %v.\n", self, m)
		r := []message.Red{message.Red{V: m[0].V}} // build r from d
		sEnd := s0.Master_1_Scatter_Red(r)
		scributil.Debugf("Worker[%d]: sent %v.\n", self, r)
		return *sEnd
	})
	wg.Done()
}
