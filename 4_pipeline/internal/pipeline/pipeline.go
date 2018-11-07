//go:generate scribblec-param.sh ../../Pipeline.scr -d ../../ -param Pipeline github.com/nickng/scribble-go-examples/4_pipeline/Pipeline -param-api W
package pipeline

import (
	"encoding/gob"
	"fmt"
	"strconv"
	"sync"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"
	"github.com/rhu1/scribble-go-runtime/test/util"

	"github.com/nickng/scribble-go-examples/4_pipeline/messages"
	"github.com/nickng/scribble-go-examples/4_pipeline/Pipeline/Pipeline"
	Head "github.com/nickng/scribble-go-examples/4_pipeline/Pipeline/Pipeline/family_1/W_1toKsub1_not_2toK"
	Middle "github.com/nickng/scribble-go-examples/4_pipeline/Pipeline/Pipeline/family_1/W_1toKsub1and2toK"
	Tail "github.com/nickng/scribble-go-examples/4_pipeline/Pipeline/Pipeline/family_1/W_2toK_not_1toKsub1"
	"github.com/nickng/scribble-go-examples/scributil"
)

var _ = shm.Dial
var _ = tcp.Dial


//*
var LISTEN = tcp.Listen
var DIAL = tcp.Dial
var FORMATTER = func() *session2.GobFormatter { return new(session2.GobFormatter) } 
/*/
var LISTEN = shm.Listen
var DIAL = shm.Dial
var FORMATTER = func() *session2.PassByPointer { return new(session2.PassByPointer) } 
//*/



const PORT = 8888

func init() {
	var foo messages.Foo
	gob.Register(&foo)
}



func Server_tail(wg *sync.WaitGroup, K int, self int) *Tail.End {
	P1 := Pipeline.New()
	R := P1.New_family_1_W_2toK_not_1toKsub1(K, self)
	var ss transport2.ScribListener
	var err error
	if ss, err = LISTEN(PORT+self); err != nil {
		panic(err)
	}
	defer ss.Close()
	if err = R.
			W_1toKsub1and2toK_Accept(self-1, ss, FORMATTER());
			//W_1toKsub1_not_2toK_Accept(self-1, ss, FORMATTER());  // Target variant (L/M) not constrained, but safe to use either
			err != nil {
		panic(err)
	}
	fmt.Println("Tail (" + strconv.Itoa(R.Self) + ") accepted", self-1, "on", PORT+self)
	end := R.Run(runTail)
	wg.Done()
	return &end
}

func runTail(s *Tail.Init) Tail.End {
	pay := make([]messages.Foo, 1)
	scributil.Delay(1500)
	end := s.W_selfsub1_Gather_Foo(pay)
	fmt.Println("Tail (" + strconv.Itoa(s.Ept.Self) + ") received:", pay)
	fmt.Println("Tail (" + strconv.Itoa(s.Ept.Self) + ") finished")
	return *end
}

func Server_middle(wg *sync.WaitGroup, K int, self int) *Middle.End {
	P1 := Pipeline.New()
	M := P1.New_family_1_W_1toKsub1and2toK(K, self)
	var ss transport2.ScribListener
	var err error
	if ss, err = LISTEN(PORT+self); err != nil {
		panic(err)
	}
	defer ss.Close()
	if self == 2 {
		if err = M.W_1toKsub1_not_2toK_Accept(self-1, ss, FORMATTER()); err != nil {
			panic(err)
		}
	} else {
		if err = M.W_1toKsub1and2toK_Accept(self-1, ss, FORMATTER()); err != nil {
			panic(err)
		}
	}
	fmt.Println("Middle (" + strconv.Itoa(M.Self) + ") accepted", self-1, "on", PORT+self)
	if self == K - 1 {
		if err := M.W_2toK_not_1toKsub1_Dial(self+1, util.LOCALHOST, PORT+self+1, DIAL, FORMATTER()); err != nil {
			panic(err)
		}
	} else {
		if err := M.W_1toKsub1and2toK_Dial(self+1, util.LOCALHOST, PORT+self+1, DIAL, FORMATTER()); err != nil {
			panic(err)
		}
	}
	fmt.Println("Middle (" + strconv.Itoa(M.Self) + ") connected to", self+1, "on", PORT+self+1)
	end := M.Run(runMiddle)
	wg.Done()
	return &end
}

func runMiddle(s *Middle.Init) Middle.End {
	pay := make([]messages.Foo, 1)
	s2 := s.W_selfsub1_Gather_Foo(pay)
	fmt.Println("Middle (" + strconv.Itoa(s.Ept.Self) + ") received:", pay)
	pay = []messages.Foo{messages.Foo{s.Ept.Self}}
	scributil.Delay(1500)
	end := s2.W_selfplus1_Scatter_Foo(pay)
	fmt.Println("Middle (" + strconv.Itoa(s.Ept.Self) + ") sent Foo:", pay)
	fmt.Println("Middle (" + strconv.Itoa(s.Ept.Self) + ") finished")
	return *end
}

func Client_head(wg *sync.WaitGroup, K int, self int) *Head.End {
	P1 := Pipeline.New()
	L := P1.New_family_1_W_1toKsub1_not_2toK(K, self)  // Endpoint needs n to check self
	if err := L.
			W_1toKsub1and2toK_Dial(self+1, util.LOCALHOST, PORT+self+1, DIAL, FORMATTER());
			//W_2toK_not_1toKsub1_Dial(self+1, util.LOCALHOST, PORT+self+1, DIAL, FORMATTER());  // Target variant (M/R) not constrained, but safe to use either
			err != nil {
		panic(err)
	}
	fmt.Println("Head (" + strconv.Itoa(L.Self) + ") connected to", self+1, "on", PORT+self+1)
	end := L.Run(runHead)
	wg.Done()
	return &end
}

func runHead(s *Head.Init) Head.End {
	pay := []messages.Foo{messages.Foo{s.Ept.Self}}
	scributil.Delay(1500)
	end := s.W_selfplus1_Scatter_Foo(pay)
	fmt.Println("Head (" + strconv.Itoa(s.Ept.Self) + ") sent Foo:", pay)
	fmt.Println("Head (" + strconv.Itoa(s.Ept.Self) + ") finished")
	return *end
}
