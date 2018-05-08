//go:generate scribblec-param.sh ../OneToMany.scr -d ../ -param Scatter -param-api Scatter A

package main

import (
	"encoding/gob"
	"log"

	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter"
	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Scatter/A_1to1"
	"github.com/nickng/scribble-go-examples/1_one-to-many/onetomany"
	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"
)

const (
	k = 2
)

func init() {
	var data onetomany.Data
	gob.Register(data)
}

func main() {
	s := Scatter.New()
	A := s.New_A_1to1(k, 1)
	for i := 0; i < k; i++ {
		if err := A.B_1toK_Dial(i+1, "localhost", 3333+i, tcp.Dial, new(session2.GobFormatter)); err != nil {
			log.Fatal(err)
		}
	}
	A.Run(func(s *A_1to1.Init_4) A_1to1.End {
		var d []onetomany.Data
		for i := 0; i < k; i++ {
			d = append(d, onetomany.Data(i))
		}
		end := s.B_1toK_Scatter_Data(d)
		return *end
	})
}
