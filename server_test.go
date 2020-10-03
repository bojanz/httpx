// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package httpx_test

import (
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"os"
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

func Test_NewServerTLS(t *testing.T) {
	var cert tls.Certificate
	server := httpx.NewServerTLS("", cert, nil)
	if server.Addr != ":https" {
		t.Errorf("got %v, want :https", server.Addr)
	}
	if server.TLSConfig == nil {
		t.Errorf("TLSConfig is not set")
	}
	if len(server.TLSConfig.Certificates) == 0 {
		t.Errorf("Certificates are not set")
	}
}

func TestServer_IsTLS(t *testing.T) {
	server := httpx.NewServer("", nil)
	if server.IsTLS() {
		t.Errorf("got true, want false")
	}

	var cert tls.Certificate
	server = httpx.NewServerTLS("", cert, nil)
	if !server.IsTLS() {
		t.Errorf("got false, want true")
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

	// Name of a non existing file
	fh, err := ioutil.TempFile("", "server-*.unix")
	if err != nil {
		t.Skip(err)
	}
	fh.Close()
	// Race condition: if the file is created between Remove and Listen, we fail.
	os.Remove(fh.Name())
	server = httpx.NewServer("unix:"+fh.Name(), http.DefaultServeMux)
	ln, err = server.Listen()
	if _, ok := ln.(*net.UnixListener); !ok {
		t.Errorf("got %T, want *net.UnixListener", ln)
	}
	if err != nil {
		t.Errorf("unexpected error %+v", err)
	}
}
