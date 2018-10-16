//go:generate scribblec-param.sh ../../Ring.scr -d ../.. -param RingProto github.com/nickng/scribble-go-examples/5_ring/Ring -param-api W

package ring

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"
	"github.com/rhu1/scribble-go-runtime/test/util"

	"github.com/nickng/scribble-go-examples/5_ring/Ring/RingProto"
	W1 "github.com/nickng/scribble-go-examples/5_ring/Ring/RingProto/family_2/W_1to1_not_2to2and2toKsub1and3toKandKtoK"
	M "github.com/nickng/scribble-go-examples/5_ring/Ring/RingProto/family_2/W_2toKsub1and3toK_not_1to1and2to2andKtoK"
	WK "github.com/nickng/scribble-go-examples/5_ring/Ring/RingProto/family_2/W_3toKandKtoK_not_1to1and2to2and2toKsub1"
	"github.com/nickng/scribble-go-examples/5_ring/messages"
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
	var bar messages.Bar
	gob.Register(&bar)
}

// self == K
func Ring_last(wg *sync.WaitGroup, K int, self int) *WK.End {
	P1 := RingProto.New()
	WK := P1.New_family_2_W_3toKandKtoK_not_1to1and2to2and2toKsub1(K, self)
	var ss transport2.ScribListener
	var err error
	if ss, err = LISTEN(PORT + self); err != nil {
		panic(err)
	}
	defer ss.Close()
	// family 1: K > 3 -- so must accept from M -- but could also use "interoperably" between families
	if err = WK.W_2toKsub1and3toK_not_1to1and2to2andKtoK_Accept(self-1, ss, FORMATTER()); err != nil {
		panic(err)
	}
	fmt.Println("WK ("+strconv.Itoa(WK.Self)+") accepted", self-1, "on", PORT+self)
	if err := WK.W_1to1_not_2to2and2toKsub1and3toKandKtoK_Dial(1, util.LOCALHOST, PORT+1, DIAL, FORMATTER()); err != nil {
		panic(err)
	}
	fmt.Println("WK ("+strconv.Itoa(WK.Self)+") connected to", 1, "on", PORT+1)
	end := WK.Run(runWK)
	wg.Done()
	return &end
}

func runWK(s *WK.Init) WK.End {
	var end *WK.End
	switch c := s.W_selfsub1_Branch().(type) {
	case *WK.Foo_W_Init:
		var x messages.Foo
		s2 := c.Recv_Foo(&x)
		fmt.Println("WK ("+strconv.Itoa(s.Ept.Self)+") received Foo:", x)
		pay := []messages.Foo{messages.Foo{X: s.Ept.Self}}
                scributil.Delay(1500)
		s = s2.W_1_Scatter_Foo(pay)
		fmt.Println("WK ("+strconv.Itoa(s.Ept.Self)+") sent Foo:", pay)
		return runWK(s)
	case *WK.Bar_W_Init:
		var x messages.Bar
		s3 := c.Recv_Bar(&x)
		fmt.Println("WK ("+strconv.Itoa(s.Ept.Self)+") received Bar:", x)
		pay := []messages.Bar{messages.Bar{Y: strconv.Itoa(s.Ept.Self)}}
                scributil.Delay(1500)
		end = s3.W_1_Scatter_Bar(pay)
		fmt.Println("WK ("+strconv.Itoa(s.Ept.Self)+") sent Foo:", pay)
                scributil.Delay(1500)
		fmt.Println("WK ("+strconv.Itoa(s.Ept.Self)+") finished")
		return *end
	default:
		log.Fatal("Shouldn't get in here: ", reflect.TypeOf(c))
	}
	fmt.Println("WK ("+strconv.Itoa(s.Ept.Self)+") finished")
	return *end
}

// K > 3
func Ring_mid(wg *sync.WaitGroup, K int, self int) *M.End {
	P1 := RingProto.New()
	M := P1.New_family_2_W_2toKsub1and3toK_not_1to1and2to2andKtoK(K, self)
	var ss transport2.ScribListener
	var err error
	if ss, err = LISTEN(PORT + self); err != nil {
		panic(err)
	}
	defer ss.Close()

if self > 2 {
if err = M.W_2toKsub1and3toK_not_1to1and2to2andKtoK_Accept(self-1, ss, FORMATTER()); err != nil { // FIXME: shouldn't have
panic(err)
}
} else {
if err = M.W_1to1_not_2to2and2toKsub1and3toKandKtoK_Accept(self-1, ss, FORMATTER()); err != nil {
panic(err)
}
}
fmt.Println("M (" + strconv.Itoa(M.Self) + ") accepted", self-1, "on", PORT+self)

if self == K-1 {
if err := M.W_3toKandKtoK_not_1to1and2to2and2toKsub1_Dial(self+1, util.LOCALHOST, PORT+self+1, DIAL, FORMATTER()); err != nil {
panic(err)
}
} else {
if err := M.W_2toKsub1and3toK_not_1to1and2to2andKtoK_Dial(self+1, util.LOCALHOST, PORT+self+1, DIAL, FORMATTER()); err != nil {
panic(err)
}
}
	
end := M.Run(runM)
	wg.Done()
	return &end
}

func runM(s *M.Init) M.End {
	var end *M.End
	switch c := s.W_selfsub1_Branch().(type) {
	case *M.Foo_W_Init: // CHECKME: case type name vs. serverWK
		var x messages.Foo
		s2 := c.Recv_Foo(&x)
		fmt.Println("M ("+strconv.Itoa(s.Ept.Self)+") received Foo:", x)
		pay := []messages.Foo{messages.Foo{X: s.Ept.Self}}
                scributil.Delay(1500)
		s = s2.W_selfplus1_Scatter_Foo(pay)
		fmt.Println("M ("+strconv.Itoa(s.Ept.Self)+") sent Foo:", pay)
		return runM(s)
	case *M.Bar_W_Init:
		var x messages.Bar
		s3 := c.Recv_Bar(&x)
		fmt.Println("M ("+strconv.Itoa(s.Ept.Self)+") received Bar:", x)
		pay := []messages.Bar{messages.Bar{Y: strconv.Itoa(s.Ept.Self)}}
                scributil.Delay(1500)
		end = s3.W_selfplus1_Scatter_Bar(pay)
		fmt.Println("M ("+strconv.Itoa(s.Ept.Self)+") sent Foo:", pay)
                scributil.Delay(1500)
		fmt.Println("M ("+strconv.Itoa(s.Ept.Self)+") finished")
		return *end
	default:
		log.Fatal("Shouldn't get in here: ", reflect.TypeOf(c))
	}
	fmt.Println("M ("+strconv.Itoa(s.Ept.Self)+") finished")
	return *end
}

// self == 1
func Ring_ini(wg *sync.WaitGroup, K int, self int) *W1.End {
	P1 := RingProto.New()
	W1 := P1.New_family_2_W_1to1_not_2to2and2toKsub1and3toKandKtoK(K, self)
	var ss transport2.ScribListener
	var err error
	if ss, err = LISTEN(PORT + self); err != nil {
		panic(err)
	}
	defer ss.Close()
	if err := W1.W_2toKsub1and3toK_not_1to1and2to2andKtoK_Dial(self+1, util.LOCALHOST, PORT+self+1, DIAL, FORMATTER()); err != nil {
		panic(err)
	}
	// FIXME: W_1to1_not_2to2and2toKsub1and3toKandKtoK_Accept ??
	fmt.Println("W1 ("+strconv.Itoa(W1.Self)+") connected to", self+1, "on", PORT+self+1)
	if err = W1.W_3toKandKtoK_not_1to1and2to2and2toKsub1_Accept(self+K-1, ss, FORMATTER()); err != nil {
		panic(err)
	}
	fmt.Println("W1 ("+strconv.Itoa(W1.Self)+") accepted", self+K-1, "on", PORT+self)
	end := W1.Run(runW1)
	wg.Done()
	return &end
}

var seed = rand.NewSource(time.Now().UnixNano())
var rnd = rand.New(seed)
var count = 1

func runW1(s *W1.Init) W1.End {
	//var end *W1.End
	if rnd.Intn(2) < 1 {
		pay := []messages.Foo{messages.Foo{X: s.Ept.Self}}
                scributil.Delay(1500)
		s2 := s.W_2_Scatter_Foo(pay)
		fmt.Println("W1 ("+strconv.Itoa(s.Ept.Self)+") sent Foo #"+strconv.Itoa(count)+":", pay)
		s = s2.W_K_Gather_Foo(pay)
		fmt.Println("W1 ("+strconv.Itoa(s.Ept.Self)+") received:", pay)
		count = count + 1
		return runW1(s)
	} else {
		pay := []messages.Bar{messages.Bar{Y: strconv.Itoa(s.Ept.Self)}}
                scributil.Delay(1500)
		s3 := s.W_2_Scatter_Bar(pay)
		fmt.Println("W1 ("+strconv.Itoa(s.Ept.Self)+") sent Bar:", pay)
		end := s3.W_K_Gather_Bar(pay)
		fmt.Println("W1 ("+strconv.Itoa(s.Ept.Self)+") received:", pay)
                scributil.Delay(1500)
		fmt.Println("W1 ("+strconv.Itoa(s.Ept.Self)+") finished")
		return *end
	}
}
