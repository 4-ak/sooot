package security

import "testing"

func TestGenerateHash(t *testing.T) {
	id, pass, salt := "helloworld", "123456", "123"

	hashed, err := CreatePass(id, pass, salt)
	if err != nil {
		t.Errorf("can't hash password : %v", err)
	}

	if err = ComparePass(id, pass, salt, hashed); err != nil {
		t.Errorf("mismatch error : %v", err)
	}
}
