# mesh

## horizontal mesh

The horizontal wave demonstrates
`(1,1)` → `(2,1)` → ... → `(K,1)`

Run distributed version by:

    go run main.go -K 4      # W4 @ horizontal-WK
    go run main.go -K 4 -x 3 # W3 @ horizontal-Wi
    go run main.go -K 4 -x 2 # W2 @ horizontal-Wi
    go run main.go -K 4      # W1 @ horizontal-W1

## vertical mesh

The column pipe demonstrates
`(1,1)` → `(1,1)` → ... → `(1,K)`

Run distributed version by:

    go run main.go -K 4      # W4 @ vertical-WK
    go run main.go -K 4 -y 3 # W3 @ vertical-Wi
    go run main.go -K 4 -y 2 # W2 @ vertical-Wi
    go run main.go -K 4      # W1 @ vertical-W1
