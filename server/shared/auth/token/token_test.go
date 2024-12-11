package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlMrYeUrmEXP1o+xV4avJ
iPdQehqz2+3E+7E0WraeYFSlbwP9sbhN1/p+xNV8XU8MlOHofWplzhpIEQFf+Dm/
G6ihnkaJXSknpdbmDBQ+z3ggMwNrkOCJ1Bi0d5qepV5W6s25rK63SgcfDILwN2HR
WK0okXcmzh7C6Y42rr+z9VZrLZG9br/KWCHxrK/vS5d7u67JIc55RkKdEDyd1RJB
5/0RJ7JrxU/xd//FsV1ubBVezsIqyfWZFFG59J+Q8N2xm1WtBBjDtsmd4twf6klx
J1rZybYOdCsdstq2d09TIS4tNb2WoE3KUE9j3AGgKOOAdM0QVjMK/XxSsxshQw2W
OQIDAQAB
-----END PUBLIC KEY-----`

func TestVerify(t *testing.T) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		t.Fatalf("cannot parse public key: %v", err)
	}
	v := &JWTTokenVerifier{
		PublicKey: pubKey,
	}
	// 有多个case需要考虑，使用表格测试
	cases := []struct {
		name string
		tkn  string
		// 测试的时间
		now  time.Time
		want string
		// 验证在某些情况下是否出错
		wantErr bool
	}{
		{
			// 正常case
			name: "valid_token",
			tkn:  "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjc1MmIxMTVlODIwYjc3NDliOTNmYWFiIn0.PFJYAwKQ_OknD0SIsEYaFz5KbtjJTh4tdB3ldOm1cbvcmkZdd6_qwGM0fZkWp7CpmvvEcEvQfEBukOr-2SOoivG-n_KUzRpeQINflolR6HnJpDvLp9WuV2CTyPrx-jODmRVNMmb5uFpfbYrITv7VDrfUE1ThtOhPlIYgIVp0TMehhEE3JI9262S-bcmZfGCyOuFaotHnTKbo7Jvgj32pWxddLygSEL1bSm3QsLW8re_4vaMTct3pcaUYWXwUXfYKI2khJhBmZYfGMsiIHiKbnjc_BYvRv0Zps_4Npc4nPCrTOem1s2xWMvuwO2wd_jpXhDArTZuABDhCl4Ze2sjNeg",
			now:  time.Unix(1516239122, 0),
			want: "6752b115e820b7749b93faab",
			// 该case不应该出错
			wantErr: false,
		},
		{
			// 超时case
			name:    "token_expired",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjc1MmIxMTVlODIwYjc3NDliOTNmYWFiIn0.PFJYAwKQ_OknD0SIsEYaFz5KbtjJTh4tdB3ldOm1cbvcmkZdd6_qwGM0fZkWp7CpmvvEcEvQfEBukOr-2SOoivG-n_KUzRpeQINflolR6HnJpDvLp9WuV2CTyPrx-jODmRVNMmb5uFpfbYrITv7VDrfUE1ThtOhPlIYgIVp0TMehhEE3JI9262S-bcmZfGCyOuFaotHnTKbo7Jvgj32pWxddLygSEL1bSm3QsLW8re_4vaMTct3pcaUYWXwUXfYKI2khJhBmZYfGMsiIHiKbnjc_BYvRv0Zps_4Npc4nPCrTOem1s2xWMvuwO2wd_jpXhDArTZuABDhCl4Ze2sjNeg",
			now:     time.Unix(1517239122, 0),
			want:    "",
			wantErr: true,
		},
		{
			// bad token case
			name:    "bad_token",
			tkn:     "bad_token",
			now:     time.Unix(1516239122, 0),
			want:    "",
			wantErr: true,
		},
		{
			// 伪造token，使用不同的私钥计算签名
			name:    "wrong_signature",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjc1MmIxMTVlODIwYjc3NDliOTNmYWFiIn0.G40igTsXSSyDur2Qv-SJSHeDIagp31veUqGBUZy3pBXiMEmhtMIod6QfxzGTsdmqDpUzRnbW00PcpKZx4t9g6Wl-S0XXw0jMj627m22d6D_3zJvvvOVluP8MH-RnQS77tQ4KfFPzpydXnsKSRAOGXOSo7GfE3FV9haoQH8SLXTQgR4C7XKN2XjQRlLS68Ro2O3AqT3QQGLG67aw_gSq_cQsJN172_dWnpWQh_4nLdGCgRjKvZU4cMdvsYhvcVKxHVPH75c35L_2F70ov0p-pBsqUEnT1km33l2qh1EOcWY1tfDGxSu6HXCLufsbXLU9kSz9dQozi2LZULBVy4uM22g",
			now:     time.Unix(1516239122, 0),
			want:    "",
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// 固定测试时间，不依赖于外界时间
			jwt.TimeFunc = func() time.Time {
				return c.now
			}
			accountID, err := v.Verify(c.tkn)
			// 如果不期望出错，但是出错了
			if !c.wantErr && err != nil {
				t.Errorf("verification failed: %v", err)
			}
			// 如果期望出错，但是没有出错

			if c.wantErr && err == nil {
				t.Errorf("want error; got no error")
			}

			if accountID != c.want {
				t.Errorf("wrong account id. want: %q, got: %q", c.want, accountID)
			}
		})
	}
}
