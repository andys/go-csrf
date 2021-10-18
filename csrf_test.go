package csrf

import (
	"testing"
)

func TestIsTokenValid(t *testing.T) {
	Setup([]byte("x"))

	if IsTokenValid("") {
		t.Fatalf("IsTokenValid on emptry string should return false!")
	}

	if IsTokenValid("x") {
		t.Fatalf("IsTokenValid on invalid string should return false!")
	}

	token := GenerateToken()
	Setup([]byte("y"))

	if IsTokenValid(token) {
		t.Fatalf("IsTokenValid after random secret changed should return false!")
	}

	Setup([]byte("b4d86f34908e66ebbf0ab78d2b80d592"))

	token = GenerateToken()
	if !IsTokenValid(token) {
		t.Fatalf("IsTokenValid on valid token %v should return true!", token)
	}
}
