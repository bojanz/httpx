# httpx [![Build Status](https://travis-ci.org/bojanz/httpx.png?branch=master)](https://travis-ci.org/bojanz/httpx) [![Coverage Status](https://coveralls.io/repos/github/bojanz/httpx/badge.svg?branch=master)](https://coveralls.io/github/bojanz/httpx?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/bojanz/httpx)](https://goreportcard.com/report/github.com/bojanz/httpx) [![GoDoc](https://godoc.org/github.com/bojanz/httpx?status.svg)](https://godoc.org/github.com/bojanz/httpx)

Provides an extended, production-ready HTTP server.

## Features

1. Production-ready defaults (TLS, timeouts), following [Cloudflare recommendations](https://blog.cloudflare.com/exposing-go-on-the-internet/).
2. Limiter for max simultaneous connections.
3. Support for systemd sockets.

```go
    // Serve r on port 80.
    err := httpx.ListenAndServe(":80", r)
    // Serve r on a systemd socket.
    err := httpx.ListenAndServe("systemd:myapp-http.socket", r)

    // Serve up to 1000 simultaneous connections on port 80.
    server := httpx.NewServer(":80", r)
    server.MaxConnections = 1000
    err := server.ListenAndServe()
```

## Alternatives

- [ory/graceful](https://github.com/ory/graceful) provides production-ready defaults and graceful shutdown.
