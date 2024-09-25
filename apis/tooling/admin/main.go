package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	if err := GenKey(); err != nil {
		fmt.Println(err)
	}
}

func GenKey() error {

	// 生成 private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("error generating private key: %w", err)
	}

	// 生成 key 存储的私钥文件
	privateFile, err := os.Create("private.pem")
	if err != nil {
		return fmt.Errorf("error creating private file: %w", err)
	}
	defer privateFile.Close()

	// 将 privateKey 包装成 pem 编码的格式
	privateBlock := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	// 将 key 写入到 privateFile 中
	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		return fmt.Errorf("encoding to private file: %w", err)
	}

	// -------------------------------------------------------------------------
	// 生成公钥
	publicFile, err := os.Create("public.pem")
	if err != nil {
		return fmt.Errorf("creating public file: %w", err)
	}
	defer publicFile.Close()

	// 从私钥创建公钥
	publicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshaling public key: %w", err)
	}

	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKey,
	}
	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		return fmt.Errorf("encoding to public file: %w", err)
	}

	fmt.Println("private and public key files generated")
	return nil
}
