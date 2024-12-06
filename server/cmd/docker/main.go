package main

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func main() {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.43"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	// 创建一个mongo容器
	resp, err := c.ContainerCreate(ctx, &container.Config{
		Image: "mongo:4.4",
		ExposedPorts: nat.PortSet{
			"27017/tcp": {},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"27017/tcp": []nat.PortBinding{
				{
					HostIP: "127.0.0.1",
					// 端口设置为0会随机分配一个可用端口
					HostPort: "0",
				},
			},
		},
	}, nil, nil, "")

	if err != nil {
		panic(err)
	}
	// 启动容器
	err = c.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("container started")
	time.Sleep(10 * time.Second)

	// 查询随机映射选择的端口
	inspRes, err := c.ContainerInspect(ctx, resp.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("listening at %+v\n", inspRes.NetworkSettings.Ports["27017/tcp"][0])

	// 这里启动的docker供测试使用，测试完毕后直接remove（不使用stop，因为stop会保留container仍会占用资源）
	fmt.Println("container removed")
	err = c.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
	if err != nil {
		panic(err)
	}
}
