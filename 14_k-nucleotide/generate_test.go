//go:generate scribblec-param.sh KNuc.scr -d . -param Proto github.com/nickng/scribble-go-examples/14_k-nucleotide/KNuc -param-api A -param-api B -param-api S
//go:generate sh -c "go run ../tools/fasta/fasta.go -n 1000 > ./input/input_1.fasta"
//go:generate sh -c "go run ../tools/fasta/fasta.go -n 12500000 > ./input/input_2.fasta"
//go:generate sh -c "go run ../tools/fasta/fasta.go -n 25000000 > ./input/input_3.fasta"
//go:generate sh -c "go run ../tools/fasta/fasta.go -n 50000000 > ./input/input_4.fasta"

package main

// This is an empty file.
//
// The _test suffix of the file name prevents the Go compiler
// from complaining about missing main() function but allows
// the go:generate directives to be found by "go generate".
