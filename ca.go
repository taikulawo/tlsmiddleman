package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"time"
)

type Certificate struct {
	certificate *x509.Certificate
	privateKey *PrivateKey
}


type CertificateAuthority struct {
	config *RuntimeConfig
	tlsConfig *TLSConfig
	privateKey *PrivateKey
	certificate *Certificate
}


func NewCA(config *RuntimeConfig, tlsConfig *TLSConfig) *CertificateAuthority{
	ca := &CertificateAuthority{
		config: config,
		tlsConfig: tlsConfig,
	}
	panic(ca.generateSelfSignCertificate())
	return ca
}

func (this *CertificateAuthority) LoadPkFromFile(filePath string) (*PrivateKey,error){
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	rsaPair, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{
		rsaPair,
	}, nil
}

func (this *CertificateAuthority) LoadCertificateFromFile(path string) (*Certificate, error) {

}

func (this *CertificateAuthority) CreatedCertificateFor(organization, host string, isCA bool, issuer *Certificate) (*Certificate, error) {
	template := &x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(int64(time.Now().UnixNano())),
		Subject: pkix.Name{
			Organization: []string{organization},
			CommonName: host,
		},
		DNSNames: []string{host},
		NotBefore: time.Now().AddDate(0,-1,0),
		NotAfter: time.Now().Add(time.Hour * 24 * 365),
		KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}
	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = []net.IP{ip}
	}
	if (isCA) {
		template.IsCA = true
		template.KeyUsage = template.KeyUsage | x509.KeyUsageCertSign
	}
	isSelfSign := issuer == nil
	if isSelfSign {
		template.ExtKeyUsage = [] x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth | x509.ExtKeyUsageClientAuth}
	}
}

// 如果自签CA证书，pkOfClient == pk
// 如果为client签发那么 pkOfClient是重新生成的密钥对，pk是CA的密钥对
func (this *CertificateAuthority) doCreateCertificate(template, issuer *x509.Certificate, pkOfClient *PrivateKey,pk *PrivateKey) (*Certificate, error){
	certBytes, err := x509.CreateCertificate(rand.Reader,template,issuer,pkOfClient.rsa.Public(),pk.rsa)
	if err != nil {
		return nil, err
	}
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, err
	}
	return &Certificate{
		certificate: cert,
		privateKey: pkOfClient,
	}, nil
}

func (this *CertificateAuthority) createTemplateFor(organization string, dnsName string) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(int64(time.Now().UnixNano())),
		Subject: pkix.Name{
			Organization: []string{organization},
			CommonName: dnsName,
		},
		DNSNames: []string{dnsName},
		NotBefore: time.Now().AddDate(0,-1,0),
		NotAfter: time.Now().Add(time.Hour * 24 * 365),
		KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}
}

func (this *CertificateAuthority) generateSelfSignCertificate() error {
	var err error
	if this.privateKey, err = this.LoadPkFromFile(this.tlsConfig.CAPrivateKeyFilePath); err != nil {
		if this.privateKey, err = this.GeneratePrivateKey(this.tlsConfig.KeyLen); err != nil {
			panic(err)
		}
		return this.privateKey.WriteToFile(this.tlsConfig.CAPrivateKeyFilePath)
	}
	if this.certificate, err = this.LoadCertificateFromFile(this.tlsConfig.CACertificateFilePath); err != nil {
		if this.certificate, err = this.CreatedCertificateFor(this.tlsConfig.CommonName); err != nil {

		}
	}
}

func (this *CertificateAuthority) GeneratePrivateKey(keyLen int) (*PrivateKey, error) {
	pair, err := rsa.GenerateKey(rand.Reader,keyLen)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{
		rsa: pair,
	}, nil
}
