package decoder

import "crypto/tls"

// RuntimeConfig is a wrapper for CA
type RuntimeConfig struct {
	Port string
}

type TLSConfig struct {
	RsaKeyPair            *PrivateKey // CA key pair
	CAPrivateKeyFilePath  string
	CACertificateFilePath string
	Organization          string
	CommonName            string // also used in DNS name
	ServerTLSConfig       *tls.Config
	KeyLen                int
}

func NewDefaultServerTlsConfig() *tls.Config {
	return &tls.Config{
		CipherSuites: []uint16{
			// 这里的顺序是我抓了个client hello包对照着写的
			// 真是长。。
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		PreferServerCipherSuites: true,
		//InsecureSkipVerify: true,
	}
}

// NewDefaultTLSConfig return configuration used in TLS handshake
func NewDefaultTLSConfig() *TLSConfig {
	return &TLSConfig{
		RsaKeyPair:            nil,
		CAPrivateKeyFilePath:  "private.key",
		CACertificateFilePath: "x509cert.crt",
		Organization:          "bytejump",
		CommonName:            "www.bytejump.com",
		ServerTLSConfig:       NewDefaultServerTlsConfig(),
		KeyLen: 4096,
	}
}
