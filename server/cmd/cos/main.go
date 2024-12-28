package main

import (
	"context"
	blobpb "coolcar/blob/api/gen/v1"
	"fmt"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8083", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	c := blobpb.NewBlobServiceClient(conn)
	ctx := context.Background()

	// res, err := c.CreateBlob(ctx, &blobpb.CreateBlobRequest{
	// 	AccountId:           "account_1",
	// 	UploadUrlTimeoutSec: 1000,
	// })

	// res, err := c.GetBlob(ctx, &blobpb.GetBlobRequest{
	// 	Id: "676f88c84aea956a1e0ed05b",
	// })

	res, err := c.GetBlobURL(ctx, &blobpb.GetBlobURLRequest{
		Id:         "676f88c84aea956a1e0ed05b",
		TimeoutSec: 1000,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
}
