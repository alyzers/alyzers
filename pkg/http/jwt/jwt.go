package jwt

import (
	"errors"
	"github.com/alyzers/alyzers/pkg/http"
	"github.com/alyzers/alyzers/pkg/log"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2023/11/14 22:40
 * @file: jwt.go
 * @description:
 */

type AuthClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

func (a *AuthClaims) Valid() error {
	return nil
}

var (
	issUser = "alyzers"
)

func keyFunc(auth http.Auth) (interface{}, error) {
	return auth.SecretKey, nil
}

// GenToken 生成 access_token 和 refresh_token
func GenToken(userId string, secretKey []byte, accessExpired, refreshExpired time.Duration) (aToken, rToken string, err error) {

	// aToken
	aClaims := &AuthClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issUser, // 签发人
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExpired * time.Minute)),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	aToken, aErr := jwt.NewWithClaims(jwt.SigningMethodHS256, aClaims).SignedString(secretKey)
	if aErr != nil {
		log.Errorf("jwt.NewWithClaims err: %v", aErr)
		return "", "", aErr
	}

	// rToken
	rClaims := jwt.RegisteredClaims{
		Issuer:    issUser,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExpired * time.Minute)),
	}
	rToken, rErr := jwt.NewWithClaims(jwt.SigningMethodHS256, rClaims).SignedString(secretKey)
	if rErr != nil {
		log.Debugf("jwt.NewWithClaims err: %v", rErr)
		return "", "", rErr
	}

	return aToken, rToken, nil
}

// ParseToken 校验 access_token
func ParseToken(aToken, secretKey string) (claims *AuthClaims, err error) {
	var token *jwt.Token
	claims = new(AuthClaims)
	token, err = jwt.ParseWithClaims(aToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid { // token 是否有效
		err = errors.New("invalid token")
	}
	return claims, nil
}

// RefreshToken 刷新 access_token
func RefreshToken(auth *http.Auth, userId, rToken string) (map[string]string, error) {
	newToken := make(map[string]string)

	// 解析刷新令牌
	var refreshClaims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(rToken, &refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(auth.SecretKey), nil
	})
	if err != nil {
		log.Errorf("jwt.ParseWithClaims err: %v", err)
		return newToken, errors.New(http.InvalidToken.Msg)
	}

	// 检查刷新令牌是否有效且未过期
	if refreshClaims.ExpiresAt == nil || time.Now().After(refreshClaims.ExpiresAt.Time) {
		log.Errorf("jwt.ParseWithClaims err: %v", err)
		return newToken, errors.New(http.InvalidToken.Msg)
	}

	// 生成新地访问令牌和刷新令牌
	newAToken, newRToken, err := GenToken(userId, []byte(auth.SecretKey), auth.AccessExpire, auth.RefreshExpire)
	if err != nil {
		return newToken, err
	}

	newToken["accessToken"] = newAToken
	newToken["refreshToken"] = newRToken

	return newToken, nil
}
