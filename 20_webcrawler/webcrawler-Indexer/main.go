package main

import (
	"sync"

	"github.com/nickng/scribble-go-examples/22_webcrawler/WebCrawler/Crawler"
	"github.com/nickng/scribble-go-examples/22_webcrawler/internal/webcrawler"
	"github.com/nickng/scribble-go-examples/scributil"
)

func main() {
	protocol := Crawler.New()
	conn, N := scributil.ParseFlags()

	// This should be connected 2/4.

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go webcrawler.Indexer(protocol, N, 1, conn, conn.BasePort(), conn, "localhost", conn.BasePort(), wg)
	wg.Wait()
}
