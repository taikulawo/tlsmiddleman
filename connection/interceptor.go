package connection

import (
	"errors"
	"fmt"
	"github.com/iamwwc/tlsmiddleman/decoder"
	"net"
	"net/http"
)

type Interceptor struct {
	config    *decoder.RuntimeConfig
	tlsConfig *decoder.TLSConfig
	CA        *decoder.CertificateAuthority
}

func (this *Interceptor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		this.PerformTLSHandshake(w, r)
		return
	}
	go NewConnectionHandler(w, r, this, nil)
}

func (this *Interceptor) PerformTLSHandshake(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	conn, err := this.Hijacker(w)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// TLS初始化完就交给Handler进行下面的握手
	go NewConnectionHandler(w, r, this, conn).TLSHandshake()
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

func NewInterceptor(c *decoder.RuntimeConfig, tlsConfig *decoder.TLSConfig) (i *Interceptor) {
	ca := decoder.NewCA(c, tlsConfig)
	return &Interceptor{
		config:    c,
		tlsConfig: tlsConfig,
		CA:        ca,
	}
}
