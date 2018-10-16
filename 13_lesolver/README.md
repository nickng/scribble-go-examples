# solver

Running the distributed version (9 process, `K=3`) in this order:

    go run main.go -K 3 -y 3      # W[3,3],W[3,1] @ solver-W1K
    go run main.go -K 3 -y 2      # W[3,2],W[3,1] @ solver-W1K
    go run main.go -K 3 -y 1      # W[3,1],W[3,2] @ solver-W1K

    go run main.go -K 3 -x 2 -y 3 # W[2,3] @ solver-Wi
    go run main.go -K 3 -x 2 -y 2 # W[2,2] @ solver-Wi
    go run main.go -K 3 -x 2 -y 1 # W[2,1] @ solver-Wi
