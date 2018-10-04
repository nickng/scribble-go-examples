//go:generate scribblec-param.sh -d ../ ../PGet.scr -param Basic github.com/nickng/scribble-go-examples/9_pget/PGet -param-api M -param-api F -param-api S
//go:generate scribblec-param.sh -d ../ ../PGet.scr -param Sync github.com/nickng/scribble-go-examples/9_pget/PGet -param-api A -param-api B

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/nickng/scribble-go-examples/9_pget/PGet/Basic"
	"github.com/nickng/scribble-go-examples/9_pget/PGet/Basic/F_1to1and1toK"
	"github.com/nickng/scribble-go-examples/9_pget/PGet/Basic/F_1toK_not_1to1"
	"github.com/nickng/scribble-go-examples/9_pget/PGet/Basic/family_1/M_1to1"
	"github.com/nickng/scribble-go-examples/9_pget/PGet/Sync"
	"github.com/nickng/scribble-go-examples/9_pget/PGet/Sync/A_1to1"
	"github.com/nickng/scribble-go-examples/9_pget/PGet/Sync/B_1to1"
	"github.com/nickng/scribble-go-examples/9_pget/http"
	"github.com/nickng/scribble-go-examples/9_pget/msg"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/shm"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"
)

var (
	// K is the number of fetchers.
	K int
	// URL is the URL to fetch.
	URL string
)

func init() {
	flag.IntVar(&K, "K", 2, "Specify number of fetchers")
	log.SetPrefix("pget: ")
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "pget [-K fetchers] URL\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
}

func parseURL(URL string) (host string, port int) {
	u, err := url.Parse(URL)
	if err != nil {
		log.Fatalf("invalid URL: %v", err)
	}
	if host = u.Hostname(); host == "" {
		log.Fatalf("invalid host: %s", URL)
	}
	if port, err = strconv.Atoi(u.Port()); err != nil {
		port, err = net.LookupPort("tcp", u.Scheme)
		if err != nil {
			log.Fatal(err)
		}
	}
	return host, port
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
	}
	URL = flag.Arg(0)
	httpHost, httpPort := parseURL(URL)

	p := Basic.New()
	M := p.New_family_1_M_1to1(K, 1) // Master M[1]

	F1 := p.New_F_1to1and1toK(K, 1)                        // Fetcher F[1]
	F2toK := make([]*F_1toK_not_1to1.F_1toK_not_1to1, K-1) // Fetcher F[2,K]
	for i := 2; i <= K; i++ {
		F2toK[i-2] = p.New_F_1toK_not_1to1(K, i)
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
	go initM(M, wg)
	go initF1(F1, httpHost, httpPort, FListeners[0], wg)
	for i := 2; i <= K; i++ {
		go initF2toK(F2toK[i-2], i, httpHost, httpPort, FListeners[i-1], wg)
	}
	wg.Wait()
}

func initM(M1 *M_1to1.M_1to1, wg *sync.WaitGroup) {
	debugf("M[1]: Connecting with F[1,K].\n")
	defer wg.Done()
	if err := M1.F_1to1and1toK_Dial(1, "inmem", 0, shm.Dial, new(session2.PassByPointer)); err != nil {
		log.Fatalf("cannot dial from M to F[1]: %v", err)
	}
	for i := 2; i <= K; i++ {
		port := i - 1
		if err := M1.F_1toK_not_1to1_Dial(i, "inmem", port, shm.Dial, new(session2.PassByPointer)); err != nil {
			log.Fatalf("cannot dial from M to F[%d]: %v", i, err)
		}
	}
	debugf("M[1]: Ready.\n")
	M1.Run(M)
	M1.Close()
}

func initF1(F *F_1to1and1toK.F_1to1and1toK, shost string, sport int, mln transport2.ScribListener, wg *sync.WaitGroup) {
	debugf("F[1]: Connecting with S[1] and M[1].\n")
	defer wg.Done()
	if err := F.S_1to1_Dial(1, shost, sport, tcp.Dial, new(http.Formatter)); err != nil {
		log.Fatalf("cannot dial from F to S: %v", err)
	}
	if err := F.M_1to1_Accept(1, mln, new(session2.PassByPointer)); err != nil {
		log.Fatalf("cannot accept connection from M to F: %v", err)
	}
	debugf("F[1]: Ready.\n")
	F.Run(F1)
	F.Close()
}

func initF2toK(F *F_1toK_not_1to1.F_1toK_not_1to1, Fid int, shost string, sport int, mln transport2.ScribListener, wg *sync.WaitGroup) {
	debugf("F[%d]: Connecting with S[1] and M[1].", Fid)
	defer wg.Done()
	if err := F.M_1to1_Accept(1, mln, new(session2.PassByPointer)); err != nil {
		log.Fatalf("cannot accept connection from M to F[%d]: %v", Fid, err)
	}
	if err := F.S_1to1_Dial(1, shost, sport, tcp.Dial, new(http.Formatter)); err != nil {
		log.Fatalf("cannot dial from F[%d] to S: %v", Fid, err)
	}
	debugf("F[%d]: Ready\n", Fid)
	F.Run(F2toK(Fid))
	F.Close()
}

// connectAB establishes a shared memory connection between AB.
func connectAB(A *A_1to1.A_1to1, B *B_1to1.B_1to1, Fid int) {
	debugf("Fetcher[%d]: Connecting Sync.A[1] and Sync.B[1].\n", Fid)
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
	debugf("F[%d]: Connected.\n", Fid)
	ln.Close()
}

// M implements Master
func M(s *M_1to1.Init) M_1to1.End {
	meta := allocMeta()
	s0 := s.F_1_Gather_Meta(meta)

	jobs := makeJobs(meta, K)
	data := allocData(K)
	s1 := s0.F_1toK_Scatter_Job(jobs).F_1toK_Gather_Data(data)

	for i := range data {
		debugf("Received from F[%d] %d bytes\n", i, len(data[i].Data))
		debugf("(%d)--------------------", i)
		fmt.Print(string(data[i].Data))
		debugf("--------------------(%d)", i)
	}
	fmt.Println()

	As := make([]*A_1to1.Init, K)
	s2 := s1.F_1toK_Gather_Sync(As)

	done := makeDone()
	for _, A := range As {
		A.B_1_Scatter_Done(done) // Notify each B of done.
	}
	return *s2
}

// F1 implements Fetcher 1
func F1(s *F_1to1and1toK.Init) F_1to1and1toK.End {
	// Make HTTP HEAD request.
	headReq := makeHeadReq(URL)
	headRes := allocResponse()
	s0 := s.S_1_Scatter_Head(headReq).S_1_Gather_Res(headRes)

	debugf("F[1]: size=%d", extractSize(headRes))

	// Notify Master of metadata.
	meta := makeMeta(URL, extractSize(headRes))
	job := allocJob()
	s1 := s0.M_1_Scatter_Meta(meta).M_1_Gather_Job(job)

	// Make HTTP GET request.
	getReq := makeGetReq(job)
	getRes := allocResponse()
	s2 := s1.S_1_Scatter_Get(getReq).S_1_Gather_Res(getRes)

	// Send data to Master.
	data := collectData(getRes)
	s3 := s2.M_1_Scatter_Data(data)

	syn := Sync.New()
	A1, B1 := syn.New_A_1to1(1), syn.New_B_1to1(1)
	connectAB(A1, B1, 1)
	s4 := s3.M_1_Scatter_Sync([]*A_1to1.Init{A1.Init()})

	B1.Run(B)
	return *s4
}

// F2toK implements Fetcher 2 ... Fetcher K
func F2toK(id int) func(s *F_1toK_not_1to1.Init) F_1toK_not_1to1.End {
	return func(s *F_1toK_not_1to1.Init) F_1toK_not_1to1.End {
		job := allocJob()
		s0 := s.M_1_Gather_Job(job)

		getReq := makeGetReq(job)
		getRes := allocResponse()
		s1 := s0.S_1_Scatter_Get(getReq).S_1_Gather_Res(getRes)

		// Send data to Master
		data := collectData(getRes)
		s2 := s1.M_1_Scatter_Data(data)

		syn := Sync.New()
		A1, B1 := syn.New_A_1to1(1), syn.New_B_1to1(1)
		connectAB(A1, B1, id)
		s3 := s2.M_1_Scatter_Sync([]*A_1to1.Init{A1.Init()})

		B1.Run(B)
		return *s3
	}
}

// B implements Sync.B as a simple wait signal.
func B(s *B_1to1.Init) B_1to1.End {
	done := allocDone()
	s0 := s.A_1_Gather_Done(done)
	return *s0
}

// ---- Helper functions ----

var debug = os.Getenv("DEBUG") == "1"

func debugf(format string, args ...interface{}) {
	if debug {
		log.Printf(format, args...)
	}
}

func allocResponse() []http.Response {
	return make([]http.Response, 1)
}

func makeHeadReq(url string) []http.HeadReq {
	return []http.HeadReq{http.Head(url)}
}

func extractSize(res []http.Response) int {
	size, err := strconv.Atoi(res[0].Header.Get("Content-Length"))
	if err != nil {
		return 0
	}
	return size
}

func makeMeta(url string, size int) []msg.Meta {
	return []msg.Meta{msg.Meta{URL: url, Size: size}}
}

func allocMeta() []msg.Meta {
	return make([]msg.Meta, 1)
}

func makeJobs(meta []msg.Meta, K int) []msg.Job {
	jobs := make([]msg.Job, K)
	fragSize := meta[0].Size / K
	for i := 0; i < K; i++ {
		jobs[i].URL = meta[0].URL
		jobs[i].RangeFrom = i * fragSize
		if i < K {
			jobs[i].RangeTo = (i+1)*fragSize - 1
		} else {
			jobs[i].RangeTo = meta[0].Size
		}
	}
	return jobs
}

func allocJob() []msg.Job {
	return make([]msg.Job, 1)
}

func makeGetReq(job []msg.Job) []http.GetReq {
	return []http.GetReq{http.Get(job[0].URL, job[0].RangeFrom, job[0].RangeTo)}
}

func collectData(res []http.Response) []msg.Data {
	return []msg.Data{msg.Data{Data: res[0].Body}}
}

func allocData(K int) []msg.Data {
	return make([]msg.Data, K)
}

func makeDone() []msg.Done {
	return []msg.Done{msg.Done{}}
}

func allocDone() []msg.Done {
	return make([]msg.Done, 1)
}
