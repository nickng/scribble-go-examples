//go:generate scribblec-param.sh ../OneToMany.scr -d ../ -param Scatter -param-api github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany B

package main

import (
	"encoding/gob"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter"
	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter/B_1toK"
	"github.com/nickng/scribble-go-examples/1_one-to-many/onetomany"
	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"
)

const k = 2

func init() {
	var data onetomany.Data
	gob.Register(data)
}

func main() {
	s := Scatter.New()
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go gather(s, 1, wg)
	go gather(s, 2, wg)
	wg.Wait()
}

func gather(s *Scatter.Scatter, id int, wg *sync.WaitGroup) {
	ln, err := tcp.Listen(3333 + id - 1)
	if err != nil {
		log.Fatalf("Cannot listen: %v", err)
	}
	B := s.New_B_1toK(k, id)
	if err := B.A_1to1_Accept(1, ln, new(session2.GobFormatter)); err != nil {
		log.Fatal(err)
	}
	B.Run(func(s *B_1toK.Init_8) B_1toK.End {
		d := make([]onetomany.Data, 1)
		end := s.A_1to1_Gather_Data(d)
		return *end
	})
	wg.Done()
}
