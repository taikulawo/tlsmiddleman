package main

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
)

const (
	PEM_HEADER_OF_PRIVATE_KEY = "RSA PRIVATE KEY"
)

type PrivateKey struct {
	rsa          *rsa.PrivateKey
}

func (this *PrivateKey) PemDecoded() []byte {
	return pem.EncodeToMemory(this.pemBlock())
}

func (this *PrivateKey) encodeToBytes(w io.Writer)  {
	pem.Encode(bufio.NewWriter(w),this.pemBlock())
}

func (this *PrivateKey) pemBlock() *pem.Block{
	return &pem.Block{
		Type:  PEM_HEADER_OF_PRIVATE_KEY,
		Bytes: x509.MarshalPKCS1PrivateKey(this.rsa),
	}
}

func (this *PrivateKey) WriteToFile(path string) error {
	file, err := os.OpenFile(path,os.O_WRONLY|os.O_CREATE|os.O_TRUNC,0600)
	if err != nil {
		return fmt.Errorf("Open %s failed. Caused By %s",path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Close file %s failed. Caused By", path, err)
		}
	}()
	return pem.Encode(file,this.pemBlock())
}