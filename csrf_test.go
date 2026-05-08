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
// called with a bare integer (24*60*60), which was interpreted as nanoseconds.
// The kataras/jwt library treats any MaxAge <= 1 second as NoMaxAge and strips
// the expiry claim entirely, so we must verify the token actually carries an
// "exp" roughly 24 hours in the future rather than just checking it parses.
func TestTokenMaxAgeIsHours(t *testing.T) {
	Setup([]byte("b4d86f34908e66ebbf0ab78d2b80d592"))

	before := time.Now().Unix()
	token := GenerateToken()

	verified, err := jwt.Verify(jwt.HS512, secret, []byte(token))
	if err != nil {
		t.Fatalf("unexpected verify error: %v", err)
	}

	exp := verified.StandardClaims.Expiry
	if exp == 0 {
		t.Fatalf("token has no expiry set (MaxAge was probably <= 1s and became NoMaxAge)")
	}

	ttl := exp - before
	const day = int64(24 * 60 * 60)
	// Allow a few seconds of slack for clock drift / test execution time.
	if ttl < day-5 || ttl > day+5 {
		t.Fatalf("expected token TTL ~24h (%d seconds), got %d seconds", day, ttl)
	}
}

// TestExpiredTokenIsInvalid directly signs a token that has already expired
// and asserts IsTokenValid rejects it, confirming expiry is actually enforced.
func TestExpiredTokenIsInvalid(t *testing.T) {
	Setup([]byte("b4d86f34908e66ebbf0ab78d2b80d592"))

	// Sign a token whose "exp" is in the past. Note: we can't use MaxAge here
	// because jwt.MaxAge silently becomes NoMaxAge for durations <= 1s.
	now := time.Now().Unix()
	claims := jwt.Claims{IssuedAt: now - 3600, Expiry: now - 60}
	raw, err := jwt.Sign(jwt.HS512, secret, claims)
	if err != nil {
		t.Fatalf("unexpected sign error: %v", err)
	}

	if IsTokenValid(string(raw)) {
		t.Fatalf("expired token should not be valid")
	}
}
