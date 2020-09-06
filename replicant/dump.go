package replicant

import (
	"bufio"
	"fmt"
	"github.com/iamwwc/tlsmiddleman/common"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/http/httputil"
)

// Dump也可以使用DumpRequest将Request对象转换成[]byte
// TLS握手结束就http.Server来处理裸HTTP
// 利用现成的*http.Request和http.ResponseWriter配置DumpResponse，DumpRequest搞
// https://golang.org/pkg/net/http/httputil/#DumpRequest
func Dump() (requestChan chan []byte, responseChan chan []byte) {
	requestChan = make(chan []byte)
	reqC := make(chan *http.Request, 2)
	go func() {
		for {
			reader := common.NewReaderHelper(requestChan)
			req, err := http.ReadRequest(bufio.NewReader(reader))
			reqC <- req
			if err != nil {
				return
			}
			fmt.Printf("Request: Header: %s\n",req.Header)
		}
	}()
	responseChan = make(chan []byte)
	go func() {
		for {
			reader := common.NewReaderHelper(responseChan)
			req := <- reqC
			resp, err := http.ReadResponse(bufio.NewReader(reader),req)
			if err != nil {
				logrus.Errorln(err)
				return
			}
			fmt.Printf("Response: Header: %s",resp.Header)
		}
	}()
	return
}

func DumpRequest(r *http.Request) ([]byte,error) {
	return httputil.DumpRequest(r, true)
}

func NewResponseFrom(conn net.Conn, r *http.Request) (*http.Response,error) {
	return http.ReadResponse(bufio.NewReader(conn),r)
}

func DumpResponse(w *http.Response) ([]byte, error) {
	return httputil.DumpResponse(w,true)
}
