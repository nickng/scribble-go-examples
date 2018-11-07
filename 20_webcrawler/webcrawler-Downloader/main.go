package main

import (
	"sync"

	"github.com/nickng/scribble-go-examples/20_webcrawler/WebCrawler/Crawler"
	"github.com/nickng/scribble-go-examples/20_webcrawler/internal/webcrawler"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	protocol := Crawler.New()
	conn, N := scributil.ParseFlags()

	// This should be connected 4/4.

	wg := new(sync.WaitGroup)
	wg.Add(N)
	for n := 1; n <= N; n++ {
		go webcrawler.DownloaderI(protocol, N, n, conn, "localhost", conn.Port(N+n), wg)
	}
	wg.Wait()
}
