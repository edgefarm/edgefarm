package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func NewRSA(privatePath, publicPath string) (string, string, error) {
	// Generate a new RSA private key with 2048 bits
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// Encode the private key to the PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyFile, err := os.Create(privatePath)
	if err != nil {
		return "", "", err
	}
	pem.Encode(privateKeyFile, privateKeyPEM)
	privateKeyFile.Close()

	// Extract the public key from the private key
	publicKey := &privateKey.PublicKey

	// Encode the public key to the PEM format
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	}
	publicKeyFile, err := os.Create(publicPath)
	if err != nil {
		fmt.Println("Error creating public key file:", err)
		os.Exit(1)
	}
	pem.Encode(publicKeyFile, publicKeyPEM)
	publicKeyFile.Close()

	priv, err := os.ReadFile(privatePath)
	if err != nil {
		return "", "", err
	}
	pub, err := os.ReadFile(publicPath)
	if err != nil {
		return "", "", err
	}
	return string(priv), string(pub), nil
}
