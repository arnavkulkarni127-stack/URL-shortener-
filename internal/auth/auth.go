package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil

}

func VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) // compares the byte of stored hash and provided password
	return err == nil
}

func GenerateJWT(userID int, secret string) (string, error) { // returns the token as string and the userId here is the logged in user, secret is the private key server uses  to sign the token
	claims := jwt.MapClaims{ //claims = payload or data
		"user_id": userID,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString, secret string) (int, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,     // string of jwt
		jwt.MapClaims{}, // empty claims object
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // is the tokens signing method HMAC
				return nil, errors.New("unexpected signing method") // reject
			}
			return []byte(secret), nil // if the signature method is right return secret as byte because this verifies that the tokenn is valid annd untamperred
		})
	if err != nil { // not parsed successfully

		return 0, err
	}
	if !token.Valid { // token.Valid = not expired, sign matched, issued-at is not in the future
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("missing user id")
	}

	return int(userID), nil
}
