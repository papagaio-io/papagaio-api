package common

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type TokenSigningData struct {
	Duration   time.Duration
	Method     jwt.SigningMethod
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	Key        []byte
}

//Token utilizzato come campo "state" nella request al authorization server
func GenerateOauth2JWTToken(sd *TokenSigningData, remoteSourceName string) (string, error) {
	return GenerateGenericJWTToken(sd, jwt.MapClaims{
		"exp":             time.Now().Add(sd.Duration).Unix(),
		"git_source_name": remoteSourceName,
	})
}

//Token fornito al front-end come risposta nella callback(login)
func GenerateLoginJWTToken(sd *TokenSigningData, userID uint64) (string, error) {
	return GenerateGenericJWTToken(sd, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(sd.Duration).Unix(),
	})
}

func GenerateGenericJWTToken(sd *TokenSigningData, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(sd.Method, claims)

	var key interface{}
	switch sd.Method {
	case jwt.SigningMethodRS256:
		key = sd.PrivateKey
	case jwt.SigningMethodHS256:
		key = sd.Key
	default:
		return "", errors.Errorf("unsupported signing method %q", sd.Method.Alg())
	}
	// Sign and get the complete encoded token as a string
	return token.SignedString(key)
}

const (
	expireTimeRange time.Duration = 5 * time.Minute
)

func IsAccessTokenExpired(expiresAt time.Time) bool {
	if expiresAt.IsZero() {
		return false
	}
	return expiresAt.Add(-expireTimeRange).Before(time.Now())
}

func ParseToken(sd *TokenSigningData, stringToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method != sd.Method {
			return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		var key interface{}
		switch sd.Method {
		case jwt.SigningMethodRS256:
			key = sd.PrivateKey
		case jwt.SigningMethodHS256:
			key = sd.Key
		default:
			return nil, errors.Errorf("unsupported signing method %q", sd.Method.Alg())
		}
		return key, nil
	})

	return token, err
}
