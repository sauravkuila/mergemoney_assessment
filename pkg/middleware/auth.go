package middleware

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(secret string, customClaims map[string]interface{}, exp time.Time) (string, error) {
	claims := jwt.MapClaims{}
	for k, v := range customClaims {
		claims[k] = v
	}
	claims["exp"] = exp.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT validates a JWT token using the provided secret.
// It returns the JWT body (claims as map[string]interface{}) and an error.
// Possible errors: "signature invalid", "token expired", "token invalid".
// If the token is expired, it returns the claims and "token expired" error.
func ValidateJWT(tokenString string, secret string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signature invalid")
		}
		return []byte(secret), nil
	})

	if err != nil {
		// Signature error
		if err.Error() == jwt.ErrSignatureInvalid.Error() {
			return nil, fmt.Errorf("signature invalid")
		}
		// return nil, fmt.Errorf("token invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("token invalid")
	}

	// Expiry validation
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return claims, fmt.Errorf("token expired")
		}
	}

	return claims, nil
}
