//go:generate scribblec-param.sh Regex.scr -d . -param Proto github.com/nickng/scribble-go-examples/15_regex-redux/Regex -param-api A -param-api B -param-api C
//go:generate sh -c "go run ../tools/fasta/fasta.go -n 1000 > ./input/input_1.fasta"
//go:generate sh -c "go run ../tools/fasta/fasta.go -n 1000000 > ./input/input_2.fasta"
//go:generate sh -c "go run ../tools/fasta/fasta.go -n 2500000 > ./input/input_3.fasta"
//go:generate sh -c "go run ../tools/fasta/fasta.go -n 5000000 > ./input/input_4.fasta"

package main

// This is an empty file.
//
// The _test suffix of the file name prevents the Go compiler
// from complaining about missing main() function but allows
// the go:generate directives to be found by "go generate".
