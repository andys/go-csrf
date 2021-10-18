package csrf

import (
	"github.com/kataras/jwt"
)

var secret []byte

func Setup(randomSecret []byte) {
	secret = randomSecret
}

func GenerateToken() string {
	token, err := jwt.Sign(jwt.HS512, secret, "", jwt.MaxAge(24*60*60))
	if err != nil {
		panic(err)
	}

	return string(token)
}

func IsTokenValid(token string) bool {
	_, err := jwt.Verify(jwt.HS512, secret, []byte(token), jwt.Plain)
	if err != nil {
		//fmt.Printf("csrf.IsTokenValid: error: %+v\n", err)
		return false
	}
	return true
}
