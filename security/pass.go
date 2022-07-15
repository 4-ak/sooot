package security

import "golang.org/x/crypto/bcrypt"

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
