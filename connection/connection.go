package connection

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/iamwwc/tlsmiddleman/common"
	"github.com/iamwwc/tlsmiddleman/decoder"
	"github.com/iamwwc/tlsmiddleman/replicant"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"
	"time"
)

func NewConnectionHandler(w http.ResponseWriter, r *http.Request, interceptor *Interceptor, conn net.Conn) *Handler {
	return &Handler{
		interceptor,
		conn,
		nil,
		w,
		r,
		false,
	}
}

type Handler struct {
	interceptor *Interceptor
	conn        net.Conn
	remote      net.Conn
	response    http.ResponseWriter
	request     *http.Request
	isHttps     bool
}

func (this *Handler) TLSHandshake() {
	this.isHttps = true
	tlsConfig := decoder.NewDefaultServerTlsConfig()
	cert, err := this.interceptor.CA.Sign(strings.Split(this.request.Host, ":")[0])
	if err != nil {
		logrus.Errorln(err)
		return
	}
	tlsCert, err := this.interceptor.CA.ToTLSCertificate(cert)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	tlsConfig.Certificates = []tls.Certificate{tlsCert}
	tlsConn := tls.Server(this.conn, tlsConfig)
	this.conn = tlsConn

	if err := tlsConn.Handshake(); err != nil {
		logrus.Errorln(err)
		this.Destroy()
		return
	}
	go this.Pipe()
}

// Pipe处理裸HTTP，并向remote转发数据
// 这里需要区分 src 是HTTPS还是HTTP，对于HTTPS需要作为TLS client和remote连接
// 如果是HTTP直接转发就行
// ResponseWriter 和 Request 都放在 this 上
func (this *Handler) Pipe() {
	remote := <-this.connectToRemote()
	if remote == nil {
		this.conn.Close()
		logrus.Debugln("Connect to remote failed, return")
		return
	}
	this.remote = remote
	this.HttpAndHttpsPipe()
	// 针对数据流的Pipe
	//this.StreamPipe()
}

func (this *Handler) HttpAndHttpsPipe() {
	if this.isHttps {
		httpHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			this.HTTPPipe(writer, request)
		})
		if err := http.Serve(this, httpHandler); err != nil {
			logrus.Debugln(err)
		}
		return
	}
	this.HTTPPipe(this.response, this.request)
}

func (this *Handler) HTTPPipe(w http.ResponseWriter, r *http.Request) {
	// dump http request and http response
	go func() {
		reqDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			logrus.Debugln(err)
		}
		if _, err := this.remote.Write(reqDump); err != nil {
			return
		}
	}()
	// 我看了一下ReadResponse的源代码。
	// 虽然参数要求Request，但如果我们只是为了要Response的值的话Request设置为nil就可以
	respFromRemote, err := http.ReadResponse(bufio.NewReader(this.remote), r)
	if err != nil {
		return
	}
	respDumped, err := httputil.DumpResponse(respFromRemote, true)
	reqConn, err := this.interceptor.Hijacker(w)
	if reqConn == nil{
		this.Destroy()
		return
	}
	if _, err := reqConn.Write(respDumped);err != nil {
		this.Destroy()
		return
	}
	go func() {
		buffer := bytes.Buffer{}
		buffer.WriteString(fmt.Sprintf("Request-Host:%s\n", r.Host))
		buffer.WriteString(fmt.Sprintf("Response-Headers:\n"))
		buffer.WriteString(fmt.Sprintf("--------------------------\n\n"))
		for k, v := range respFromRemote.Header {
			buffer.WriteString(fmt.Sprintf("%s:%s\n", k, v))
		}
		fmt.Fprint(os.Stdout, buffer.String())
	}()
}
func (this *Handler) Destroy() {
	if this.conn != nil {
		this.conn.Close()
	}
	if this.remote != nil {
		this.remote.Close()
	}
	// close tls server
	this.Close()
}
func (this *Handler) StreamPipe() {
	chan1 := common.ChannelFromConn(this.conn)
	chan2 := common.ChannelFromConn(this.remote)
	reqChan, respChan := replicant.Dump()
	defer func() {
		this.remote.Close()
		this.conn.Close()
		reqChan <- nil
		respChan <- nil
	}()
	for {
		select {
		case b1 := <-chan1:
			if b1 == nil {
				return
			}
			respChan <- b1[:]
			this.remote.Write(b1)
		case b2 := <-chan2:
			if b2 == nil {
				return
			}
			reqChan <- b2[:]
			this.conn.Write(b2)
		}
	}
}

func (this Handler) connectToRemote() <-chan net.Conn {
	c := make(chan net.Conn, 1)
	go func() {
		var conn net.Conn
		var err error
		target := this.request.URL.Host
		matched, _ := regexp.MatchString(":[0-9]+$", target)
		if this.isHttps {
			if !matched {
				target += ":443"
			}
			conn, err = tls.Dial("tcp", target, decoder.NewDefaultServerTlsConfig())
			if err != nil {
				logrus.Errorln(err)
				c <- nil
				return
			}
		} else {
			if !matched {
				target += ":80"
			}
			conn, err = net.DialTimeout("tcp", target, time.Second*60)
			if err != nil {
				logrus.Errorln(err)
				c <- nil
				return
			}
		}
		c <- conn
	}()
	return c
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
	return nil
}

func (this *Handler) Addr() net.Addr {
	return nil
}
