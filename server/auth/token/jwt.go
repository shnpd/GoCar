package token

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTTokenGen struct {
	privatekey *rsa.PrivateKey
	issuer     string
	// 获取当前时间时，如果直接使用time.Now()那么在测试时每次的时间都是不一样的，我们不希望依赖于外界的时间，因此这里使用了一个函数来分别对测试和正式环境进行处理
	nowFunc func() time.Time
}

// 构造函数
func NewJWTTokenGen(issuer string, privateKey *rsa.PrivateKey) *JWTTokenGen {
	return &JWTTokenGen{
		privatekey: privateKey,
		issuer:     issuer,
		nowFunc:    time.Now,
	}
}

func (t *JWTTokenGen) GenerateToken(accountID string, expire time.Duration) (string, error) {
	nowSec := t.nowFunc().Unix()
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		Issuer:    t.issuer,
		IssuedAt:  nowSec,
		ExpiresAt: nowSec + int64(expire.Seconds()),
		Subject:   accountID,
	})

	// 签名
	return tkn.SignedString(t.privatekey)
}
