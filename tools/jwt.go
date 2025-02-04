package tools

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type CustomClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func NewAccessToken() (string, error) {
	claims := &CustomClaims{
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "raythx98@gmail.com",
			Subject:   "user",
			Audience:  []string{"raythx98@gmail.com"},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(1 * time.Minute)},
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	hmacSecret := []byte("secret")
	tokenString, _ := token.SignedString(hmacSecret)
	// TODO: handle error
	//if err != nil {
	//	errorMsg := "NewCustomToken: failed to create jwt token"
	//	s.BaseService.Logger.Error(errorMsg, zap.String("userId", userId), zap.Error(err))
	//	return "", errors.Wrap(err, errorMsg)
	//}

	return tokenString, nil
}

func NewRefreshToken() (string, error) {
	claims := &CustomClaims{
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "raythx98@gmail.com",
			Subject:   "user",
			Audience:  []string{"raythx98@gmail.com"},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(10 * time.Minute)},
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	hmacSecret := []byte("secret")
	tokenString, _ := token.SignedString(hmacSecret)
	// TODO: handle error
	//if err != nil {
	//	errorMsg := "NewCustomToken: failed to create jwt token"
	//	s.BaseService.Logger.Error(errorMsg, zap.String("userId", userId), zap.Error(err))
	//	return "", errors.Wrap(err, errorMsg)
	//}

	return tokenString, nil
}

func IsValidAccessToken(bearerAuthToken string) error {
	hmacSecret := []byte("secret")

	// Parse and verify the JWT token using the secret key and multiple issuers from _commonSecrets
	var token *jwt.Token
	token, err := jwt.ParseWithClaims(
		bearerAuthToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return hmacSecret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer("raythx98@gmail.com"))
	if err != nil {
		//s.BaseService.Logger.Debug("filterValidBearerAuthTokens: invalid token",
		//	zap.Error(err), zap.String("token", bearerAuthToken))
		return fmt.Errorf("cannot parse token: %v", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid && claims.TokenType == "access" {
		return nil
	} else {
		return fmt.Errorf("invalid token")
	}
}

func IsValidRefreshToken(bearerAuthToken string) error {
	hmacSecret := []byte("secret")

	// Parse and verify the JWT token using the secret key and multiple issuers from _commonSecrets
	var token *jwt.Token
	token, err := jwt.ParseWithClaims(
		bearerAuthToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return hmacSecret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer("raythx98@gmail.com"))
	if err != nil {
		//s.BaseService.Logger.Debug("filterValidBearerAuthTokens: invalid token",
		//	zap.Error(err), zap.String("token", bearerAuthToken))
		return fmt.Errorf("cannot parse token: %v", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid && claims.TokenType == "refresh" {
		return nil
	} else {
		return fmt.Errorf("invalid token")
	}
}
