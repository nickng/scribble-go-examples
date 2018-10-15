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

## scatter mesh

A scatter style mesh `M[1,1]` → `W[K,K]`
(note that by default it is a square mesh)

Run distributed version by:

    go run main.go -K 2           # M      @ scatter-M
    go run main.go -K 2 -x 2 -y 2 # W(2,2) @ scatter-W
    go run main.go -K 2 -x 2 -y 1 # W(2,1) @ scatter-W
    go run main.go -K 2 -x 1 -y 2 # W(1,2) @ scatter-W
    go run main.go -K 2 -x 1 -y 1 # W(1,1) @ scatter-W

## gather mesh

A gather style mesh `W[K,K]` → `M[1,1]`
(note that by default it is a square mesh)

Run distributed version by:

    go run main.go -K 2           # M      @ gather-M
    go run main.go -K 2 -x 2 -y 2 # W(2,2) @ gather-W
    go run main.go -K 2 -x 2 -y 1 # W(2,1) @ gather-W
    go run main.go -K 2 -x 1 -y 2 # W(1,2) @ gather-W
    go run main.go -K 2 -x 1 -y 1 # W(1,1) @ gather-W

# diagonal mesh

A diagonal mesh pattern

`(1,1)` → `(2,2)` → ... → `(K,K)`

Run distributed version by (note `-xy` specifies both x and y coordinate):

    go run main.go -K 4       # W4 @ diag-WKK
    go run main.go -K 4 -xy 3 # W3 @ diag-Wii
    go run main.go -K 4 -xy 2 # W2 @ diag-Wii
    go run main.go -K 4       # W1 @ diag-W11
