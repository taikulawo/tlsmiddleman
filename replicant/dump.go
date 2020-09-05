package replicant

import (
	"bufio"
	"fmt"
	"github.com/iamwwc/tlsmiddleman/common"
	"net/http"
)

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
				return
			}
			fmt.Printf("Response: Header: %s",resp.Header)
		}
	}()
	return
}
