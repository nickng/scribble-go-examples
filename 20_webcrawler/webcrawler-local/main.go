package main

import (
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/20_webcrawler/WebCrawler/Crawler"
	"github.com/nickng/scribble-go-examples/20_webcrawler/internal/webcrawler"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	protocol := Crawler.New()
	conn, N := scributil.ParseFlags()

	// Connections:
	//
	// offset 0: Indexer[1] -> Master[1]
	// offset 1..N: n of Parser[I] -> Indexer[1]
	// offset N+1..2N: n of Downloader[I] -> Parser[I] with matching I

	wg := new(sync.WaitGroup)
	wg.Add(1 + 1 + N + N)
	go webcrawler.Master(protocol, N, 1, conn, conn.BasePort(), wg)
	go func() {
		// Phase 1 - connect indexer to Master
		time.Sleep(100 * time.Millisecond)
		webcrawler.Indexer(protocol, N, 1, conn, conn.BasePort(), conn, "localhost", conn.BasePort(), wg)
	}()
	for n := 1; n <= N; n++ {
		go func(n int) {
			// Phase 2 - connect Parsers to Indexer.
			time.Sleep(200 * time.Millisecond)
			webcrawler.ParserI(protocol, N, n, conn, conn.Port(N+n), conn, "localhost", conn.Port(n), wg)
		}(n)
	}
	for n := 1; n <= N; n++ {
		go func(n int) {
			// Phase 3 - connect Downloader to Parser
			time.Sleep(300 * time.Millisecond)
			webcrawler.DownloaderI(protocol, N, n, conn, "localhost", conn.Port(N+n), wg)
		}(n)
	}
	wg.Wait()
}
