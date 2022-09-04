package security

import (
	"crypto/rand"
	"crypto/rsa"
	"strings"

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

func CreatePass(plains ...string) ([]byte, error) {

	return bcrypt.GenerateFromPassword([]byte(strings.Join(plains, "")), 11)
}

func ComparePass(hashed []byte, plains ...string) error {
	return bcrypt.CompareHashAndPassword(hashed, []byte(strings.Join(plains, "")))
}
