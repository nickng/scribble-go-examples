package foreach

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/9_pget/PGet/Foreach"
	"github.com/nickng/scribble-go-examples/9_pget/PGet/Foreach/family_1/F_1to1and1toK"
	"github.com/nickng/scribble-go-examples/9_pget/PGet/Foreach/family_1/F_1toK_not_1to1"
	"github.com/nickng/scribble-go-examples/9_pget/PGet/Foreach/family_1/M_1to1"
	"github.com/nickng/scribble-go-examples/9_pget/PGet/Sync"
	"github.com/nickng/scribble-go-examples/9_pget/PGet/Sync/A_1to1"
	"github.com/nickng/scribble-go-examples/9_pget/PGet/Sync/B_1to1"
	"github.com/nickng/scribble-go-examples/9_pget/http"
	"github.com/nickng/scribble-go-examples/9_pget/internal/pget"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"
)

var _ = time.After

var (
	// URL is the URL to fetch.
	URLclient string
)

func RunClient(K int, URL string) {
	URLclient = URL

	httpHost, httpPort := pget.ParseURL(URL)

	p := Foreach.New()
	M := p.New_family_1_M_1to1(K, 1) // Master M[1]

	F1 := p.New_family_1_F_1to1and1toK(K, 1)               // Fetcher F[1]
	F2toK := make([]*F_1toK_not_1to1.F_1toK_not_1to1, K-1) // Fetcher F[2,K]
	for i := 2; i <= K; i++ {
		F2toK[i-2] = p.New_family_1_F_1toK_not_1to1(K, i)
	}
	FListeners := make([]transport2.ScribListener, K)
	var err error
	for i := range FListeners {
		FListeners[i], err = shm.Listen(i)
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
	}

	wg := new(sync.WaitGroup)
	wg.Add(1 + K)
	go initM(M, wg, K)
	go initF1(F1, httpHost, httpPort, FListeners[0], wg)
	for i := 2; i <= K; i++ {
		go initF2toK(F2toK[i-2], i, httpHost, httpPort, FListeners[i-1], wg)
	}
	wg.Wait()
}

func initM(M1 *M_1to1.M_1to1, wg *sync.WaitGroup, K int) {
	pget.Debugf("M[1]: Connecting with F[1,K].\n")
	defer wg.Done()
	if err := M1.F_1to1and1toK_Dial(1, "inmem", 0, shm.Dial, new(session2.PassByPointer)); err != nil {
		log.Fatalf("cannot dial from M to F[1]: %v", err)
	}
	for i := 2; i <= K; i++ {
		time.Sleep(100 * time.Millisecond)
		port := i - 1
		if err := M1.F_1toK_not_1to1_Dial(i, "inmem", port, shm.Dial, new(session2.PassByPointer)); err != nil {
			log.Fatalf("cannot dial from M to F[%d]: %v", i, err)
		}
	}
	pget.Debugf("M[1]: Ready.\n")
	M1.Run(M)
}

func initF1(F *F_1to1and1toK.F_1to1and1toK, shost string, sport int, mln transport2.ScribListener, wg *sync.WaitGroup) {
	pget.Debugf("F[1]: Connecting with S[1] and M[1].\n")
	defer wg.Done()
	if err := F.M_1to1_Accept(1, mln, new(session2.PassByPointer)); err != nil {
		log.Fatalf("cannot accept connection from M to F: %v", err)
	}
	if err := F.S_1to1_Dial(1, shost, sport, tcp.Dial, new(http.Formatter)); err != nil {
		log.Fatalf("cannot dial from F to S: %v", err)
	}
	pget.Debugf("F[1]: Ready.\n")
	F.Run(F1)
	mln.Close()
}

func initF2toK(F *F_1toK_not_1to1.F_1toK_not_1to1, Fid int, shost string, sport int, mln transport2.ScribListener, wg *sync.WaitGroup) {
	pget.Debugf("F[%d]: Connecting with S[1] and M[1].", Fid)
	defer wg.Done()
	if err := F.M_1to1_Accept(1, mln, new(session2.PassByPointer)); err != nil {
		log.Fatalf("cannot accept connection from M to F[%d]: %v", Fid, err)
	}
	if err := F.S_1to1_Dial(1, shost, sport, tcp.Dial, new(http.Formatter)); err != nil {
		log.Fatalf("cannot dial from F[%d] to S: %v", Fid, err)
	}
	pget.Debugf("F[%d]: Ready\n", Fid)
	F.Run(F2toK(Fid))
	mln.Close()
}

// connectAB establishes a shared memory connection between AB.
func connectAB(A *A_1to1.A_1to1, B *B_1to1.B_1to1, Fid int, K int) {
	pget.Debugf("Fetcher[%d]: Connecting Sync.A[1] and Sync.B[1].\n", Fid)
	ln, err := shm.Listen(K + Fid)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	ch := make(chan struct{})
	go func() {
		A.B_1to1_Accept(1, ln, new(session2.PassByPointer))
		close(ch)
	}()
	go func() {
		B.A_1to1_Dial(1, "", K+Fid, shm.Dial, new(session2.PassByPointer))
	}()
	<-ch
	pget.Debugf("F[%d]: Connected.\n", Fid)
	ln.Close()
}

// M implements Master
func M(s *M_1to1.Init) M_1to1.End {
	meta := pget.AllocMeta()
	s0 := s.F_1_Gather_Meta(meta)

	jobs := pget.MakeJobs(meta, s.Ept.K)
	data := pget.AllocData(s.Ept.K)
	s1 := s0.F_1toK_Scatter_Job(jobs).F_1toK_Gather_Data(data)

	for i := range data {
		pget.Debugf("Received from F[%d] %d bytes\n", i, len(data[i].Data))
		pget.Debugf("(%d)--------------------", i)
		fmt.Print(string(data[i].Data))
		pget.Debugf("--------------------(%d)", i)
	}
	fmt.Println()

	As := make([]*A_1to1.Init, s.Ept.K)
	s2 := s1.F_1toK_Gather_Sync(As)

	done := pget.MakeDone()
	for _, A := range As {
		A.B_1_Scatter_Done(done) // Notify each B of done.
	}
	return *s2
}

// F1 implements Fetcher 1
func F1(s *F_1to1and1toK.Init) F_1to1and1toK.End {
	// Make HTTP HEAD request.
	headReq := pget.MakeHeadReq(URLclient)
	headRes := pget.AllocResponse()
	s0 := s.S_1_Scatter_Head(headReq).S_1_Gather_Res(headRes)

	pget.Debugf("F[1]: size=%d", pget.ExtractSize(headRes))

	// Notify Master of metadata.
	meta := pget.MakeMeta(URLclient, pget.ExtractSize(headRes))
	job := pget.AllocJob()
	s1 := s0.M_1_Scatter_Meta(meta).M_1_Gather_Job(job)

	// Make HTTP GET request.
	getReq := pget.MakeGetReq(job)
	getRes := pget.AllocResponse()
	s2 := s1.S_1_Scatter_Get(getReq).S_1_Gather_Res(getRes)

	// Send data to Master.
	data := pget.CollectData(getRes)
	s3 := s2.M_1_Scatter_Data(data)

	syn := Sync.New()
	A1, B1 := syn.New_A_1to1(1), syn.New_B_1to1(1)
	connectAB(A1, B1, 1, s.Ept.K)
	s4 := s3.M_1_Scatter_Sync([]*A_1to1.Init{A1.Init()})

	B1.Run(B)
	return *s4
}

// F2toK implements Fetcher 2 ... Fetcher K
func F2toK(id int) func(s *F_1toK_not_1to1.Init) F_1toK_not_1to1.End {
	return func(s *F_1toK_not_1to1.Init) F_1toK_not_1to1.End {
		job := pget.AllocJob()
		s0 := s.M_1_Gather_Job(job)

		getReq := pget.MakeGetReq(job)
		getRes := pget.AllocResponse()
		s1 := s0.S_1_Scatter_Get(getReq).S_1_Gather_Res(getRes)

		// Send data to Master
		data := pget.CollectData(getRes)
		s2 := s1.M_1_Scatter_Data(data)

		syn := Sync.New()
		A1, B1 := syn.New_A_1to1(1), syn.New_B_1to1(1)
		connectAB(A1, B1, id, s.Ept.K)
		s3 := s2.M_1_Scatter_Sync([]*A_1to1.Init{A1.Init()})

		B1.Run(B)
		return *s3
	}
}

// B implements Sync.B as a simple wait signal.
func B(s *B_1to1.Init) B_1to1.End {
	done := pget.AllocDone()
	s0 := s.A_1_Gather_Done(done)
	return *s0
}
