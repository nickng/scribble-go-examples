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

	// This should be connected 1/4.

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go webcrawler.Master(protocol, N, 1, conn, conn.BasePort(), wg)
	wg.Wait()
}
