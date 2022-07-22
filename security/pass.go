package security

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

type RSAKey struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

var RSAkey RSAKey

func KeyGen() error {
	private, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}
	RSAkey.Private = private
	RSAkey.Public = &private.PublicKey
	return nil
}

func EncrpytionWithBase64(data []byte) (string, error) {
	r, err := rsa.EncryptPKCS1v15(rand.Reader, RSAkey.Public, data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(r), nil
}

func DecrptionWithBase64(data string) ([]byte, error) {
	r, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, RSAkey.Private, r)
}

func CreatePlainPass(id, pass, salt string) []byte {
	data := append([]byte(id), []byte(pass)...)
	return append(data, []byte(salt)...)
}

func CreatePass(id, pass, salt string) ([]byte, error) {
	return bcrypt.GenerateFromPassword(CreatePlainPass(id, pass, salt), 11)
}

func ComparePass(id, pass, salt string, hashed []byte) error {
	return bcrypt.CompareHashAndPassword(hashed, CreatePlainPass(id, pass, salt))
}
