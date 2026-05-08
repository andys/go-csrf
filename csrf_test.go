package csrf

import (
	"testing"
	"time"

	"github.com/kataras/jwt"
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

// TestTokenMaxAgeIsHours guards against the regression where jwt.MaxAge was
// called with a bare integer (24*60*60), which was interpreted as nanoseconds
// and caused generated tokens to expire almost immediately.
func TestTokenMaxAgeIsHours(t *testing.T) {
	Setup([]byte("b4d86f34908e66ebbf0ab78d2b80d592"))

	token := GenerateToken()
	// Sleep well past any sub-second expiry (e.g. the old 86.4µs bug) before
	// checking. A valid 24h token must still verify after this.
	time.Sleep(100 * time.Millisecond)

	if !IsTokenValid(token) {
		t.Fatalf("token should still be valid 100ms after generation; MaxAge is probably not in hours")
	}
}

// TestExpiredTokenIsInvalid directly signs a token that has already expired
// and asserts IsTokenValid rejects it, confirming expiry is actually enforced.
func TestExpiredTokenIsInvalid(t *testing.T) {
	Setup([]byte("b4d86f34908e66ebbf0ab78d2b80d592"))

	raw, err := jwt.Sign(jwt.HS512, secret, jwt.Map{}, jwt.MaxAge(time.Nanosecond))
	if err != nil {
		t.Fatalf("unexpected sign error: %v", err)
	}
	time.Sleep(10 * time.Millisecond)

	if IsTokenValid(string(raw)) {
		t.Fatalf("expired token should not be valid")
	}
}
