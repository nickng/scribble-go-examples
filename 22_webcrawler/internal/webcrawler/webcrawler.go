//go:generate scribblec-param.sh ../../WebCrawler.scr -d ../../ -param Crawler github.com/nickng/scribble-go-examples/22_webcrawler/WebCrawler -param-api Downloader -param-api Parser -param-api Indexer -param-api Master

// Package webcrawler provides the implementation for the WebCrawler protocol.
//
// There are a total of 4 roles:
//   Downloader - a fetcher that downloads the web page
//   Parser - a thread for each Downloader to parse the downloaded page
//   Indexer - a single coordinator for collection results from all Parsers
//   Master - a master that collects from all results from Indexer
package webcrawler

import (
	"encoding/gob"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/22_webcrawler/WebCrawler/Crawler"
	"github.com/nickng/scribble-go-examples/22_webcrawler/WebCrawler/Crawler/Downloader_1toN"
	"github.com/nickng/scribble-go-examples/22_webcrawler/WebCrawler/Crawler/Indexer_1to1"
	"github.com/nickng/scribble-go-examples/22_webcrawler/WebCrawler/Crawler/Master_1to1"
	"github.com/nickng/scribble-go-examples/22_webcrawler/WebCrawler/Crawler/Parser_1toN"
	"github.com/nickng/scribble-go-examples/22_webcrawler/message"
	"github.com/nickng/scribble-go-examples/scributil"
)

func init() {
	gob.Register(new(message.Parse))
	gob.Register(new(message.Index))
	gob.Register(new(message.URL))
}

// DownloaderI implements Downloader[I:1,N]
// DownloaderI only knows about ParserI and is client to -> ParserI
func DownloaderI(p *Crawler.Crawler, N, self int, parser scributil.ClientConn, host string, port int, wg *sync.WaitGroup) {
	DownloaderI := p.New_Downloader_1toN(N, self)

	scributil.Debugf("[connection] Downloader[%d]: Dialling to Parser[%d] (%s:%d).\n", self, self, host, port)
	if err := DownloaderI.Parser_1toN_Dial(self, host, port, parser.Dial, parser.Formatter()); err != nil {
		log.Fatalf("cannot connect: %v", err)
	}
	scributil.Debugf("[connection] Downloader[%d]: Connected.\n", self)

	DownloaderI.Run(func(s *Downloader_1toN.Init) Downloader_1toN.End {
		parsed := []message.Parse{message.Parse{}}
		sEnd := s.Parser_selfplus0_Scatter_Parse(parsed)
		scributil.Debugf("Downloader[%d]: sent %v.\n", self, parsed)
		return *sEnd
	})
	wg.Done()
}

// ParserI implements Parser[I:1,N]
// ParserI is server to <- DownloaderI
// ParserI is client to -> Indexer[1]
func ParserI(p *Crawler.Crawler, N, self int, downloader scributil.ServerConn, downloaderPort int, indexer scributil.ClientConn, indexerHost string, indexerPort int, wg *sync.WaitGroup) {
	ParserI := p.New_Parser_1toN(N, self)

	scributil.Debugf("[connection] Parser[%d]: Dialling to Indexer (%s:%d).\n", self, indexerHost, indexerPort)
	if err := ParserI.Indexer_1to1_Dial(1, indexerHost, indexerPort, indexer.Dial, indexer.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	scributil.Debugf("[connection] Parser[%d]: Listening for Downloader[%d] at :%d.\n", self, self, downloaderPort)
	ln, err := downloader.Listen(downloaderPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	if err := ParserI.Downloader_1toN_Accept(self, ln, downloader.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	scributil.Debugf("[connection] Parser[%d]: Connected.\n", self)

	ParserI.Run(func(s *Parser_1toN.Init) Parser_1toN.End {
		// allocate space to store parsed data
		parsed := make([]message.Parse, 1)
		s0 := s.Downloader_selfplus0_Gather_Parse(parsed)
		scributil.Debugf("Parser[%d]: received %v.\n", self, parsed)
		// create index
		index := []message.Index{message.Index{}}
		sEnd := s0.Indexer_1_Scatter_Index(index)
		scributil.Debugf("Parser[%d]: sent %v.\n", self, index)
		return *sEnd
	})
	wg.Done()
}

// Indexer implements Indexer[1]
// Indexer is server to <- Parser[1,N]
// Indexer is client to -> Master[1]
func Indexer(p *Crawler.Crawler, N, self int, parser scributil.ServerConn, baseport int, master scributil.ClientConn, host string, port int, wg *sync.WaitGroup) {
	Indexer := p.New_Indexer_1to1(N, self)

	scributil.Debugf("[connection] Indexer: Dialling to Master (%s:%d).\n", host, port)
	if err := Indexer.Master_1to1_Dial(1, host, port, master.Dial, master.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	wgSvr := new(sync.WaitGroup)
	wgSvr.Add(N)
	mu := new(sync.Mutex) // for concurrent access of connection map
	for n := 1; n <= N; n++ {
		go func(n int) {
			scributil.Debugf("[connection] Indexer: Listening for Parser[%d] at :%d.\n", n, baseport+n)
			ln, err := parser.Listen(baseport + n)
			if err != nil {
				log.Fatalf("cannot listen: %v", err)
			}
			mu.Lock()
			defer mu.Unlock()
			if err := Indexer.Parser_1toN_Accept(n, ln, parser.Formatter()); err != nil {
				log.Fatalf("cannot accept: %v", err)
			}
			wgSvr.Done()
		}(n)
	}
	wgSvr.Wait()
	scributil.Debugf("[connection] Indexer: Connected.\n")

	Indexer.Run(func(s *Indexer_1to1.Init) Indexer_1to1.End {
		// Allocate space to store received indices (from N Parsers).
		index := make([]message.Index, N)
		s0 := s.Parser_1toN_Gather_Index(index)
		scributil.Debugf("Indexer: received %v.\n", index)

		// Completed URL to send to Master.
		url := []message.URL{message.URL{URL: "http://example.com/"}}
		sEnd := s0.Master_1_Scatter_URL(url)
		scributil.Debugf("Indexer: sent %v.\n", url)
		return *sEnd
	})
	wg.Done()
}

// Master implements Master[1]
// Master is server to <- Indexer[1]
func Master(p *Crawler.Crawler, N, self int, indexer scributil.ServerConn, port int, wg *sync.WaitGroup) {
	Master := p.New_Master_1to1(self)

	scributil.Debugf("[connection] Master: listening for Indexer at :%d.\n", port)
	ln, err := indexer.Listen(port)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	if err := Master.Indexer_1to1_Accept(1, ln, indexer.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	scributil.Debugf("[connection] Master: Connected.\n")

	Master.Run(func(s *Master_1to1.Init) Master_1to1.End {
		// Allocate memory to receive URL from Indexer[1]
		url := make([]message.URL, 1)
		sEnd := s.Indexer_1_Gather_URL(url)
		scributil.Debugf("Master: received %v.\n", url)
		return *sEnd
	})
	wg.Done()
}
