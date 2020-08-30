package connection

import (
	"crypto/tls"
	"fmt"
	"github.com/iamwwc/tlsmiddleman/decoder"
	"io"
	"net"
	"net/http"
)

func NewConnectionHandler(w http.ResponseWriter, r *http.Request, interceptor *Interceptor, conn *net.Conn) *Handler {
	return &Handler{
		interceptor,
		*conn,
		w,
		r,
	}
}

type Handler struct {
	interceptor *Interceptor
	conn        net.Conn
	response    http.ResponseWriter
	request     *http.Request
}

func (this *Handler) TLSHandshake() {
	tlsConfig := decoder.NewDefaultServerTlsConfig()
	cert, err := this.interceptor.CA.Sign(this.request.Host)
	if err != nil {
		fmt.Println(err)
		return
	}
	tlsCert, err := this.interceptor.CA.ToTLSCertificate(cert)
	if err != nil {
		fmt.Println(err)
		return
	}
	tlsConfig.Certificates = []tls.Certificate{tlsCert}
	tlsConn := tls.Server(this.conn, tlsConfig)
	this.conn = tlsConn
	handler := http.HandlerFunc(func(resp http.ResponseWriter, r *http.Request) {
		this.Pipe()
	})
	http.Serve(this, handler)
}

// Pipe处理裸HTTP，并向remote转发数据
// 这里需要区分 src 是HTTPS还是HTTP，对于HTTPS需要作为TLS client和remote连接
// 如果是HTTP直接转发就行
// ResponseWriter 和 Request 都放在 this 上
func (this *Handler) Pipe() {

}

func (this *Handler) Accept() (net.Conn, error) {
	if this.conn != nil {
		c := this.conn
		this.conn = nil
		return c, nil
	}
	return nil, io.EOF
}

func (this *Handler) Close() error {
	return this.conn.Close()
}

func (this *Handler) Addr() net.Addr {
	return nil
}
