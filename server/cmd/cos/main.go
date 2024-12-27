package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func main() {
	u, err := url.Parse("https://coolcar-1311261643.cos.ap-nanjing.myqcloud.com")
	if err != nil {
		panic(err)
	}
	b := &cos.BaseURL{BucketURL: u}
	// 永久密钥
	secretId := "AKIDaBWbPxHK7dvxiCQ4SZQ0JL6anslEWaPz"
	secretKey := "TklXdXAYgYAXhgRxFKEsGOIwKXUBdppd"
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	})

	// 获取预签名 URL
	name := "abc.jpg"
	presignedURL, err := client.Object.GetPresignedURL(context.Background(),
		http.MethodGet,
		name, secretId, secretKey,
		1*time.Hour, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(presignedURL)
}
