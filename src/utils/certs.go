package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

// Credits: https://github.com/kgretzky/pwndrop/blob/master/core/gen_cert.go
func GenerateTLSCertificate(common string) (*tls.Certificate, error) {
	private_key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	notBefore := time.Now()
	aYear := time.Duration(10*365*24) * time.Hour
	notAfter := notBefore.Add(aYear)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	if common == "" {
		common = genRandomString(8)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:            []string{},
			Locality:           []string{},
			Organization:       []string{},
			OrganizationalUnit: []string{},
			CommonName:         common,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &private_key.PublicKey, private_key)
	if err != nil {
		return nil, err
	}

	ret_tls := &tls.Certificate{
		Certificate: [][]byte{cert},
		PrivateKey:  private_key,
	}
	return ret_tls, nil
}

func genRandomString(n int) string {
	const lb = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		t := make([]byte, 1)
		rand.Read(t)
		b[i] = lb[int(t[0])%len(lb)]
	}
	return string(b)
}
