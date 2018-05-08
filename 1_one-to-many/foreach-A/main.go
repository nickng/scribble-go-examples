//go:generate scribblec-param.sh ../OneToMany.scr -d ../ -param Foreach -param-api github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany A

package main

import (
	"encoding/gob"
	"log"

	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Foreach"
	"github.com/nickng/scribble-go-examples/1_one-to-many/OneToMany/Foreach/A_1to1"
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
	A := Foreach.New().New_A_1to1(k, 1)
	for i := 0; i < k; i++ {
		if err := A.B_1toK_Dial(i+1, "localhost", 3333+i, tcp.Dial, new(session2.GobFormatter)); err != nil {
			log.Fatal(err)
		}
	}
	A.Run(func(s *A_1to1.Init_8) A_1to1.End {
		var d []onetomany.Data
		for i := 0; i < k; i++ {
			d = append(d, onetomany.Data{V: i})
		}
		end := s.Foreach(func(s *A_1to1.Init_6) A_1to1.End {
			end := s.B_ItoI_Scatter_Data(d)
			return *end
		})
		return *end
	})
}
