package certificate

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/acme"
)

/*
Let's Encrypt 证书申请
*/

// 生成私钥
func generatePrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}
	return privateKey, nil
}

// 生成证书请求
func generateCSR(privateKey *rsa.PrivateKey, domain string) ([]byte, error) {
	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: domain,
		},
		DNSNames: []string{domain},
	}
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate request: %v", err)
	}
	return csrBytes, nil
}

// 申请证书
func requestCertificate(domain string) (*x509.Certificate, *rsa.PrivateKey, error) {
	client := &acme.Client{
		DirectoryURL: acme.LetsEncryptURL,
	}

	privateKey, err := generatePrivateKey()
	if err != nil {
		return nil, nil, err
	}

	csrBytes, err := generateCSR(privateKey, domain)
	if err != nil {
		return nil, nil, err
	}

	account := &acme.Account{
		Contact: []string{"mailto:your-email@example.com"},
	}

	ctx := context.Background()
	account, err = client.Register(ctx, account, acme.AcceptTOS)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register account: %v", err)
	}

	authz, err := client.Authorize(ctx, domain)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to authorize domain: %v", err)
	}

	const DNS01 = "dns-01"
	var dnsChallenge *acme.Challenge
	for _, challenge := range authz.Challenges {
		if challenge.Type == DNS01 {
			dnsChallenge = challenge
			break
		}
	}

	if dnsChallenge == nil {
		return nil, nil, fmt.Errorf("dns-01 challenge not found")
	}

	challenge, err := client.Accept(ctx, dnsChallenge)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to accept challenge: %v", err)
	}
	_ = challenge

	_, err = client.WaitAuthorization(ctx, authz.URI)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to wait for authorization: %v", err)
	}

	cert, _, err := client.CreateCert(ctx, csrBytes, 0, true)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %v", err)
	}

	certBlock, _ := pem.Decode(cert[0])
	if certBlock == nil {
		return nil, nil, fmt.Errorf("failed to decode certificate PEM")
	}

	certificate, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	return certificate, privateKey, nil
}

// 保存证书和私钥到文件
func saveCertificateAndKey(cert *x509.Certificate, privateKey *rsa.PrivateKey, certFile, keyFile string) error {
	certOut, err := os.Create(certFile)
	if err != nil {
		return fmt.Errorf("failed to open cert file for writing: %v", err)
	}
	defer certOut.Close()

	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	if err != nil {
		return fmt.Errorf("failed to write data to cert file: %v", err)
	}

	keyOut, err := os.Create(keyFile)
	if err != nil {
		return fmt.Errorf("failed to open key file for writing: %v", err)
	}
	defer keyOut.Close()

	err = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	if err != nil {
		return fmt.Errorf("failed to write data to key file: %v", err)
	}

	return nil
}

// 自动申请证书并保存
// err := AutoRequestCertificate("your-domain.com", "path/to/cert.pem", "path/to/key.pem")
func AutoRequestCertificate(domain, certFile, keyFile string) error {
	cert, privateKey, err := requestCertificate(domain)
	if err != nil {
		return fmt.Errorf("failed to request certificate: %v", err)
	}

	err = saveCertificateAndKey(cert, privateKey, certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to save certificate and key: %v", err)
	}

	log.Printf("Certificate and key for domain %s saved to %s and %s", domain, certFile, keyFile)
	return nil
}
