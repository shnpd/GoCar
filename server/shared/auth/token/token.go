package token

import (
	"crypto/rsa"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// JWTTokenVerifier verifies jwt access tokens.
type JWTTokenVerifier struct {
	PublicKey *rsa.PublicKey
}

// Verify verifies a token and returns account id.
func (v *JWTTokenVerifier) Verify(token string) (string, error) {
	// Parse验证token结构
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{},
		func(*jwt.Token) (interface{}, error) {
			return v.PublicKey, nil
		})
	if err != nil {
		return "", fmt.Errorf("cannot parse token: %v", err)
	}
	// 验证签名
	if !t.Valid {
		return "", fmt.Errorf("token not valid")
	}

	// 验证标准claim
	clm, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("token claim is not StandardClaims")
	}

	// 验证payload中的内容,如是否超时
	if err := clm.Valid(); err != nil {
		return "", fmt.Errorf("claim not valid: %v", err)
	}

	return clm.Subject, nil
}
