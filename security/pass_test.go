package security

import "testing"

func TestGenerateHash(t *testing.T) {
	id, pass := "helloworld", "123456"

	hashed, err := CreatePass(id, pass)
	if err != nil {
		t.Errorf("can't hash password : %v", err)
	}

	if err = ComparePass(hashed, id, pass); err != nil {
		t.Errorf("mismatch error : %v", err)
	}
}
