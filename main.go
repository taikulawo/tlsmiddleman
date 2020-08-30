package main

import (
	"github.com/iamwwc/tlsmiddleman/connection"
	"github.com/iamwwc/tlsmiddleman/decoder"
	"net/http"
)

func main() {
	c := &decoder.RuntimeConfig{
		Port: "8080",
	}
	tlsConfig := decoder.NewDefaultTLSConfig()
	interceptor := connection.NewInterceptor(c, tlsConfig)
	s := &http.Server{
		Addr:    ":" + c.Port,
		Handler: interceptor,
	}
	s.ListenAndServe()
}
