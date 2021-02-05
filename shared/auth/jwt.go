package auth

import (
	"time"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

//CustomClaims struct stores the email and other claims of the token
type CustomClaims struct {
	Email string
	jwt.StandardClaims
}

func getSecretKey() string {

	return configuration.GetResp().Token.SecretKey
}

//GenerateToken function generates a new jwt token given an email and the expiry for the token
func GenerateToken(email string, expiry time.Duration) (string, error) {
	//writing the claims part
	claims := &CustomClaims{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * expiry).Unix(),
			Issuer:    "notification-microservice",
			IssuedAt:  time.Now().Unix(),
		},
	}
	//token generated
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoding token by signing with secret key
	t, err := token.SignedString([]byte(getSecretKey()))

	return t, err
}

//ValidateToken function validates the token and returns the token details if it was valid
func ValidateToken(encodedToken string) (*jwt.Token, error) {
	//Parse decodes the token and sends it to func for
	//validation
	return jwt.ParseWithClaims(encodedToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, errors.Errorf("Invalid token", token.Header["alg"])

		}
		return []byte(getSecretKey()), nil
	})

}
