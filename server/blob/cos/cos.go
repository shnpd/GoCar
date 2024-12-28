package cos

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// Service is the COS service.
type Service struct {
	client *cos.Client
	secID  string
	secKey string
}

// NewService creates a new COS service.
func NewService(addr, secID, secKey string) (*Service, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse url: %v", err)
	}
	b := &cos.BaseURL{BucketURL: u}

	return &Service{
		client: cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  secID,
				SecretKey: secKey,
			},
		}),
		secID:  secID,
		secKey: secKey,
	}, nil
}
func (s *Service) SignURL(c context.Context, method, path string, timeout time.Duration) (string, error) {
	u, err := s.client.Object.GetPresignedURL(c, method, path, s.secID, s.secKey, timeout, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
func (s *Service) Get(c context.Context, path string) (io.ReadCloser, error) {
	res, err := s.client.Object.Get(c, path, nil)
	var b io.ReadCloser
	// 不管出不出错都将res.Body返回让外面的defer函数可以close
	if res != nil {
		b = res.Body
	}
	if err != nil {
		return b, err
	}
	if res.StatusCode >= 400 {
		return b, fmt.Errorf("got err response: %+v", res)
	}
	return b, nil
}
