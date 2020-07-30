// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package httpx_test

import (
	"net"
	"net/http"
	"testing"

	"github.com/bojanz/httpx"
)

func Test_NewServer(t *testing.T) {
	server := httpx.NewServer("", nil)
	if server.Addr != ":http" {
		t.Errorf("got %v, want :http", server.Addr)
	}
	if server.Handler != nil {
		t.Errorf("got %v, want nil", server.Handler)
	}

	server = httpx.NewServer("0.0.0.0:8080", http.DefaultServeMux)
	if server.Addr != "0.0.0.0:8080" {
		t.Errorf("got %v, want 0.0.0.0:8080", server.Addr)
	}
	if server.Handler != http.DefaultServeMux {
		t.Errorf("got %T, want *http.ServeMux", server.Handler)
	}
	if server.ReadTimeout == 0 {
		t.Errorf("ReadTimeout is not set")
	}
	if server.WriteTimeout == 0 {
		t.Errorf("WriteTimeout is not set")
	}
	if server.IdleTimeout == 0 {
		t.Errorf("IdleTimeout is not set")
	}
	if server.TLSConfig == nil {
		t.Errorf("TLSConfig is not set")
	}
}

func TestServer_Listen(t *testing.T) {
	// Invalid systemd socket.
	server := httpx.NewServer("systemd:myapp-http", nil)
	_, err := server.Listen()
	if err == nil {
		t.Errorf("got nil, want error")
	}
	want := "listen systemd myapp-http: socket not found"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}

	// Invalid TCP address.
	server = httpx.NewServer("INVALID", http.DefaultServeMux)
	_, err = server.Listen()
	if _, ok := err.(*net.OpError); !ok {
		t.Errorf("got %T, want *net.OpError", err)
	}

	// Valid TCP address.
	server = httpx.NewServer("0.0.0.0:8080", http.DefaultServeMux)
	ln, err := server.Listen()
	if _, ok := ln.(*net.TCPListener); !ok {
		t.Errorf("got %T, want *net.TCPListener", ln)
	}
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}
