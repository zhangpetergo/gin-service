package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/open-policy-agent/opa/rego"
	"os"
	"time"
)

func main() {
	if err := GenToken(); err != nil {
		fmt.Println(err)
	}
}

//go:embed rego/authentication.rego
var opaAuthentication string

func GenToken() error {

	// Generating a token requires defining a set of claims. In this applications
	// case, we only care about defining the subject and the user in question and
	// the roles they have on the database. This token will expire in a year.
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)
	claims := struct {
		jwt.RegisteredClaims
		Roles []string
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "123455674389",
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"ADMIN"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// 带上 私钥库的id
	token.Header["kid"] = "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"

	privatePem, err := os.ReadFile("zarf/keys/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1.pem")
	if err != nil {
		return fmt.Errorf("reading private pem: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePem)
	if err != nil {
		return fmt.Errorf("parsing private pem: %w", err)
	}

	// 使用私钥签名
	str, err := token.SignedString(privateKey)
	if err != nil {
		return fmt.Errorf("signing token: %w", err)
	}

	fmt.Printf("-----BEGIN TOKEN-----\n%s\n-----END TOKEN-----\n\n", str)

	// -------------------------------------------------------------------------
	// 生成该私钥的公钥，然后使用公钥进行验证 token 是否合法

	// Marshal the public key from the private key to PKIX.
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshaling public key: %w", err)
	}
	// Construct a PEM block for the public key.
	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	var b bytes.Buffer

	// Write the public key to the public key file.
	if err := pem.Encode(&b, &publicBlock); err != nil {
		return fmt.Errorf("encoding to public file: %w", err)
	}

	// -------------------------------------------------------------------------
	// 使用 rego 验证 token

	ctx := context.Background()
	query := fmt.Sprintf("x = data.%s.%s", "ardan.rego", "auth")

	fmt.Println("query:", query)

	q, err := rego.New(
		rego.Query(query),
		rego.Module("policy.rego", opaAuthentication),
	).PrepareForEval(ctx)
	if err != nil {
		return err
	}

	input := map[string]any{
		"Key":   b.String(),
		"Token": str,
		"ISS":   "service project",
	}

	results, err := q.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	if len(results) == 0 {
		return errors.New("no results")
	}

	fmt.Println("results:", results)

	result, ok := results[0].Bindings["x"].(bool)
	if !ok || !result {
		return fmt.Errorf("bindings results[%v] ok[%v]", results, ok)
	}

	fmt.Println("\nTOKEN VALIDATED!")

	return nil

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
