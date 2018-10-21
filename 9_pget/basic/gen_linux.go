// +build !windows

//go:generate scribblec-param.sh -d ../ ../PGet.scr -param Sync github.com/nickng/scribble-go-examples/9_pget/PGet -param-api A -param-api B
//go:generate scribblec-param.sh -d ../ ../PGet.scr -param Basic github.com/nickng/scribble-go-examples/9_pget/PGet -param-api M -param-api F -param-api S

package main
