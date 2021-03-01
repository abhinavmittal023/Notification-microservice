package auth

import (
	"time"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// CustomClaims struct stores the UserID and other claims of the token
type CustomClaims struct {
	UserID    uint64
	Role      int
	TokenType string
	jwt.StandardClaims
}

func getSecretKey() string {
	return configuration.GetResp().Token.SecretKey
}

// GenerateRefreshToken generates the refresh token with userID and given expiry
func GenerateRefreshToken(userID uint64, expiry time.Duration) (string, error) {
	return GenerateToken(userID, 0, constants.TokenType().Refresh, expiry)
}

// GenerateAccessToken generates the access token with userID, role and given expiry
func GenerateAccessToken(userID uint64, role int, expiry time.Duration) (string, error) {
	return GenerateToken(userID, role, constants.TokenType().Access, expiry)
}

// GenerateValidationToken generates the validation token with userID and given expiry
func GenerateValidationToken(userID uint64, expiry time.Duration) (string, error) {
	return GenerateToken(userID, 0, constants.TokenType().Validation, expiry)
}

// GenerateToken function generates a new jwt token given userID, role, tokentype and the expiry for the token
func GenerateToken(userID uint64, role int, tokenType string, expiry time.Duration) (string, error) {
	// writing the claims part
	claims := &CustomClaims{
		userID,
		role,
		tokenType,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * expiry).Unix(),
			Issuer:    "notification-microservice",
			IssuedAt:  time.Now().Unix(),
		},
	}
	// token generated
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// encoding token by signing with secret key
	t, err := token.SignedString([]byte(getSecretKey()))

	return t, err
}

// ValidateToken function validates the token and returns the token details if it was valid
func ValidateToken(encodedToken string) (*jwt.Token, error) {
	// Parse decodes the token and sends it to func for
	// validation
	return jwt.ParseWithClaims(encodedToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, errors.Errorf("Invalid token %s", token.Header["alg"])

		}
		return []byte(getSecretKey()), nil
	})
}
