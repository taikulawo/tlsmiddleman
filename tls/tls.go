package tls

import (
	"crypto/x509"
	"net/http"
)


type HTTPHandler struct{}

func (this HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		this.ServeForTLS(w,r)
		return
	}
	
}


func (this HTTPHandler) ServeForTLS(w http.ResponseWriter, r *http.Request) {
	
}

func (this HTTPHandler) CreateCertificateFor(template *x509.Certificate, issuerCert *x509.Certificate, )[]byte  {
	
}
