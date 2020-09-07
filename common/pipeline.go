package common

type Callable interface {
	Call(data[]byte, next Callable) error
}

func Compose(middlewares []Callable) func([]byte, Callable) error {
	index := 0
	return func(data []byte, next Callable) error {
		fn := middlewares[index]
		if index == len(middlewares) {
			fn = next
		}
		return fn.Call(data,middlewares[index + 1])
	}
}