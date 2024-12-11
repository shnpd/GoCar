package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEA0PbXQ+E849UM3Db5QvDmlgO9lBFT42crEv/UVxC6nU1O/LHf
JuYQQl844sRgiqkY2EoHEY5TEkCLrJCKFvHjK/rz+7KW7cWSrIi62mciBiKpN9+M
8+KPIycLLqMyAYNDLxSK4zP5oUxk9ufyR/guEGtx3egfkQSnuSuu3NOnC6/wisMh
xfHcv85UZblNjNw0d6VZCFv6OD7BA3fCzMtjWNcgaWDLhIf76LhOF/xiQVtAttXk
bmVZAq3CSZcc3uZ4gV6pVquYly31ADMh3ST0eoayeQlzQsSR/3akYSH7rVl12cbX
hShzg2qptRYpV5nlvk2zLob65+J6brLKeAeFKQIDAQABAoIBAQCnkwwWN25pFtV1
U/CYqi+AZgeF0k9/saBtYBOcrqG4u+J36vyVaTHoyAwKboktKWXvLY51mbksje50
uITE2b3f6yP12MYIb8Sr8ApIUySJ3wns8K1Md4dqHUluYRkc9XLPMp4ejfPGUkk1
Z9R3uqLWMBPkbN3DogQPDuTv2hu/1ii6LhtdJueXJ60JeSg/f6k4iTUqmIsM7onk
VvSIolu4h6sDvAvf6cPhVcYMX50NsK4QWBAqTqA7k4POF4G6GcyfpPhW6wsHtiLT
UyPmjZZiREtiQra3xM38zsP6XmunEACBOJXyo0V+acIV51UJHvUy0zsioK24nKPk
+4XL9O4xAoGBAPE75AEHFmNhzyfc5z0k5574BvjpKYKbOFKUFDNwKriUnFpeu4IB
qwMn9r2MAJJBa87ZU/hmvRYVamWh43AZPK6URYPei3Zkx4kgAJzfr0BxclTAoPEa
Q3gWZz7Kyvrf+v5KdHBKIWUV6xyvi/3xXsajRNJTOQJt+7c82lpc0M9DAoGBAN3B
SYW8D2A0DSSeTPrpfD9W4qOYiPfUi/X4nEYLrVR1OZtV6rXhmxPv/ygQMZ7L3Z4+
qeuM/fWGlXNxiizGJxO0eFNfwQrSeYrUJYwD+ssz33CHiD45mBVw0iKnflnWN2dc
2mLPviqI73HmznxLa8f76qsrU4CavWo+EJPDt6UjAoGBALcYK/ATvwxjWmX5JpGk
ByEDQ9d8sQLuaQtUVRjNk79RHHMC3/LG7VOR65bmQjC/8uGm1jL9V5sBWiYVf5i/
ru0aoMB8EwIjb8dGQPDQXiXddVeadQ1KT9c/udFQ/kr34XtMa6Amw74DqtW0794o
azT9CulQPT7410q2p0xakiodAoGBAI6cdo8GEZFaCDXLhFA9BgWp99kkCLnxPA10
0/OMIO1W8jQ+vdI4g8X/q0V1nKO2EhFp1VdB9jUHV2fF1gnTaWAXyOY9h2VRVI7+
viLckFZMRt8Umn8CBReM1ecpdm2KAVl88olndrHeHLdQquxuiKM1FIb+NxZW7Ye8
4GZXEOynAoGBAIvczZOQZr5EGS9NUgPejhgVJCCjaM+0898SbqJafVz/KQWMLtEg
LNVNZ5LG7/6xPczxvxhku/NXvgnslUx1RdnBHEku/eLajND04yrlb4aZjdsHH6H2
Sx9/DBCOxcVNaN/EYVFqIIljWamG5Wa/Q+HC8I9EcVWZZ2EtGB+DlWnp
-----END RSA PRIVATE KEY-----`

func TestGenerateToken(t *testing.T) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		// 准备工作出错使用Fatalf，测试会立即停止
		t.Fatalf("parse private key failed: %v", err)
	}
	g := NewJWTTokenGen("coolcar/auth", key)
	// 固定时间
	g.nowFunc = func() time.Time {
		return time.Unix(1516239022, 0)
	}
	tkn, err := g.GenerateToken("6752b115e820b7749b93faab", 2*time.Hour)
	if err != nil {
		// 实际测试的函数出错使用Errorf，测试还会继续执行，观察得到的测试结果
		t.Errorf("generate token failed: %v", err)
	}

	want := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjc1MmIxMTVlODIwYjc3NDliOTNmYWFiIn0.nQJwMsyO4sTsfa0KDWMzob-ApZsKb3HI7YZPVWlfw1BJJQbdMVTq7KlV529hSpW2c7Hj5cFeS6G1rQEWhR6PHfjl7mO6OBH9EfkLEwJCTsDGBdDyGZAsNisbodTGu2y8ZYfgvpFKmc3H1uR75bZsOlu0lB7BBhOsGa_yFJ0ciexczqcB-0N5f3ACYng543XSot4-01tAHezuRJBURld_D9D66mgtYUOJt451Ai3KYP0ikj5Ji69FOPtwfsKI9qTEF8TD1OVJxkb0ejx5vvulCduPUM87WWaV831D3ma7dGuAFIifCS-kpVLxGaw_3C5k1AimPtL00UpAURHWXopOcw"

	if tkn != want {
		t.Errorf("token mismatch, want: %q, got: %q", want, tkn)
	}

}
