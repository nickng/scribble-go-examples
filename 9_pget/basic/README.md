# pget

## Usage

    go generate # to generate API
    go run main.go [-K number-of-fetchers] -- URL

    DEBUG=1 go run main.go [-K number-of-fetchers] -- URL

- `K` is 2 if unspecified
- `URL` is the URL to fetch (e.g. `http://example.com/`)
- Set the `DEBUG=1` environmental variable prefix to turn on debug messages
