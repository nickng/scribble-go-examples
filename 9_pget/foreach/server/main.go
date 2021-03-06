//go:generate scribblec-param.sh -d ../../ ../../PGet.scr -param Foreach github.com/nickng/scribble-go-examples/9_pget/PGet -param-api M -param-api F -param-api S

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nickng/scribble-go-examples/9_pget/internal/foreach"
)

var (
	// K is the number of fetchers.
	K int
	// URL is the URL to fetch. -- N.B. for server, this is only used to set the port
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

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
	}
	URL = flag.Arg(0)

	foreach.RunServer(K, URL)
}
