package common

type BytePool struct {
	w int
	wcap int
	c chan []byte
}

func (this *BytePool) Get() []byte {
	select {
	
	}
}
