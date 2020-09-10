package main

import (
	"github.com/iamwwc/tlsmiddleman/connection"
	"github.com/iamwwc/tlsmiddleman/decoder"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: false,
	})
	c := &decoder.RuntimeConfig{
		Port: "8000",
	}
	tlsConfig := decoder.NewDefaultTLSConfig()
	interceptor := connection.NewInterceptor(c, tlsConfig)
	s := &http.Server{
		Addr:    "localhost:" + c.Port,
		Handler: interceptor,
	}
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
