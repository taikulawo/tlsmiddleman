package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)


type Interceptor struct {
	 config    *RuntimeConfig
	 tlsConfig *TLSConfig
	 CA        *CertificateAuthority

}

func (this *Interceptor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		this.PerformTLSHandshake(w,r)
		return
	}

}

func (this *Interceptor) PerformTLSHandshake(w http.ResponseWriter, r *http.Request)  {
	w.WriteHeader(http.StatusOK)
	conn, err := this.Hijacker(w)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// TLS初始化完就交给Handler进行下面的握手
	NewConnectionHandler(this, conn).TLSHandshake()
}

func (this Interceptor) Hijacker(w http.ResponseWriter) (*net.Conn, error) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return nil, errors.New("unable to hijacker connection from HTTP")
	}
	conn, _, err := hijacker.Hijack()
	if err != nil {
		return nil, fmt.Errorf("failed to take over the TCP connection from hihacker %s", err)
	}
	return &conn, nil
}


func (this *Interceptor) HttpProxy(w http.ResponseWriter, r *http.Request) {

}

func NewInterceptor(c *RuntimeConfig, tlsConfig *TLSConfig) (i *Interceptor){
	ca := NewCA(c, tlsConfig)
	return &Interceptor{
		config:    c,
		tlsConfig: tlsConfig,
		CA: ca,
	}

}

func main() {
	c := &RuntimeConfig{
		Port: "8080",
	}
	tlsConfig := NewDefaultTLSConfig()
	interceptor := NewInterceptor(c,tlsConfig)
	s := &http.Server{
		Addr: ":" + c.Port,
		Handler: interceptor,
	}
	s.ListenAndServe()
}