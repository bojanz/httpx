# httpx [![Build Status](https://travis-ci.org/bojanz/httpx.png?branch=master)](https://travis-ci.org/bojanz/httpx) [![Go Report Card](https://goreportcard.com/badge/github.com/bojanz/httpx)](https://goreportcard.com/report/github.com/bojanz/httpx) [![GoDoc](https://godoc.org/github.com/bojanz/httpx?status.svg)](https://godoc.org/github.com/bojanz/httpx)

Provides an extended, production-ready HTTP server.

## Features

1. Production-ready defaults (TLS, timeouts), following [Cloudflare recommendations](https://blog.cloudflare.com/exposing-go-on-the-internet/).
2. Limiter for max simultaneous connections.
3. Support for systemd sockets.

```go
    // Serve r on port 80.
    server := httpx.NewServer(":80", r)
    err := server.ListenAndServe()
    // Serve r on a systemd socket (FileDescriptorName=myapp-http).
    server := httpx.NewServer("systemd:myapp-http", r)
    err := server.ListenAndServe()

    // Serve r on port 443.
    server := httpx.NewServer(":443", r)
    err := server.ListenAndServeTLS("/srv/cert.pem", "/srv/key.pem")
    // Serve r on a systemd TLS socket (FileDescriptorName=myapp-https).
    server := httpx.NewServer("systemd:myapp-https", r)
    err := server.ListenAndServeTLS("/srv/cert.pem", "/srv/key.pem")

    // Serve up to 1000 simultaneous connections on port 8080.
    server := httpx.NewServer(":8080", r)
    server.MaxConnections = 1000
    err := server.ListenAndServe()
```

## Systemd setup

/etc/systemd/system/myapp.service:
```
[Unit]
Description=MyApp
Requires=myapp-http.socket myapp-https.socket

[Service]
Type=simple
ExecStart=/usr/bin/myapp serve
NonBlocking=true
Restart=always

[Install]
WantedBy=multi-user.target
```

/etc/systemd/system/myapp-http.socket:
```
[Unit]
Description=MyApp HTTP socket
PartOf=myapp.service

[Socket]
ListenStream=80
NoDelay=true
Service=myapp.service
FileDescriptorName=myapp-http

[Install]
WantedBy=sockets.target
```

/etc/systemd/system/myapp-https.socket:
```
[Unit]
Description=MyApp HTTPS socket
PartOf=myapp.service

[Socket]
ListenStream=443
NoDelay=true
Service=myapp.service
FileDescriptorName=myapp-https

[Install]
WantedBy=sockets.target
```

Additional resources:
- https://www.darkcoding.net/software/systemd-socket-activation-in-go/
- https://vincent.bernat.ch/en/blog/2018-systemd-golang-socket-activation

## Alternatives

- [ory/graceful](https://github.com/ory/graceful) provides production-ready defaults and graceful shutdown.
