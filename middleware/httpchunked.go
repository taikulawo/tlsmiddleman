package middleware

import "github.com/iamwwc/tlsmiddleman/common"

// 用struct的好处是可以可以存储状态
type HttpChunkedMiddleware struct {
	store []byte
}

func (this *HttpChunkedMiddleware) Call(data []byte, next common.Callable) error {
	//next.Call(data, next)
	return nil
}
