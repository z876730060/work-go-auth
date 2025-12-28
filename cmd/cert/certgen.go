package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {
	// 解析命令行参数
	domain := flag.String("domain", "www.website.cn", "Certificate domain name")
	flag.Parse()

	// 生成ECDSA私钥
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// 设置证书模板
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:  []string{"Website Org"},
			CommonName:    *domain, // 网站域名
			Country:       []string{"CN"},
			Province:      []string{"Beijing"},
			Locality:      []string{"Beijing"},
			StreetAddress: []string{"Website Street"},
			PostalCode:    []string{"100000"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // 有效期1年
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, // 服务器认证
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback}, // 允许的IP地址
		DNSNames:              []string{*domain, "website.cn"},                 // 允许的域名
	}

	// 自签名证书
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	// 保存证书到文件
	certOut, err := os.Create("server.crt")
	if err != nil {
		log.Fatalf("Failed to open server.crt for writing: %v", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		log.Fatalf("Failed to write cert to file: %v", err)
	}
	if err := certOut.Close(); err != nil {
		log.Fatalf("Failed to close server.crt: %v", err)
	}
	log.Printf("Certificate saved to server.crt")

	// 保存私钥到文件
	keyDER, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Failed to marshal private key: %v", err)
	}
	keyOut, err := os.Create("server.key")
	if err != nil {
		log.Fatalf("Failed to open server.key for writing: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}); err != nil {
		log.Fatalf("Failed to write key to file: %v", err)
	}
	if err := keyOut.Close(); err != nil {
		log.Fatalf("Failed to close server.key: %v", err)
	}
	log.Printf("Private key saved to server.key")

	// 保存公钥到文件
	pubKeyDer, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatalf("Failed to marshal public key: %v", err)
	}
	pubOut, err := os.Create("server.pub")
	if err != nil {
		log.Fatalf("Failed to open server.pub for writing: %v", err)
	}
	if err := pem.Encode(pubOut, &pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyDer}); err != nil {
		log.Fatalf("Failed to write public key to file: %v", err)
	}
	if err := pubOut.Close(); err != nil {
		log.Fatalf("Failed to close server.pub: %v", err)
	}
	log.Printf("Public key saved to server.pub")
}