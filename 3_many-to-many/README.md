### Note

When running the program, parameter `K` is specified by the flag `-K number`.
Since there are two parameters in this protocol, `K` is split into `M` and `N`
by evenly splitting `K` (e.g. `K=4`, `M=2 N=2`), but if K is odd `M+N` will
still be guaranteed to be equal to `K`.
