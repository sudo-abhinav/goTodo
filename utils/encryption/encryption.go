package encryption

import (
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

}

// GenerateJWT creates a new JWT token with claims
func GenerateJWT(userID, email, sessionID string) (string, error) {
	claims := jwt.MapClaims{
		"userID":    userID,
		"email":     email,
		"sessionID": sessionID,
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(time.Minute * 20).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

//func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
//	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorMalformed)
//		}
//		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
//	})
//
//	if err != nil {
//		var vErr *jwt.ValidationError
//		if errors.As(err, &vErr) {
//			switch {
//			case vErr.Errors&jwt.ValidationErrorExpired != 0:
//				return nil, fmt.Errorf("token has expired")
//			case vErr.Errors&jwt.ValidationErrorNotValidYet != 0:
//				return nil, fmt.Errorf("token is not valid yet")
//			default:
//				return nil, fmt.Errorf("invalid token")
//			}
//		}
//		return nil, fmt.Errorf("could not parse token: %w", err)
//	}
//
//	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//		return claims, nil
//	}
//
//	return nil, fmt.Errorf("invalid token claims")
//}
