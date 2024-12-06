package mongotesting

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	image         = "mongo:4.4"
	containerPort = "27017/tcp"
)

// RunWithMongoInDocker runs the tests with a mongodb instance in a docker container.
func RunWithMongoInDocker(m *testing.M, mongoURI *string) int {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.43"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	// 创建一个mongo容器
	resp, err := c.ContainerCreate(ctx, &container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			containerPort: {},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{
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
	// 在函数退出时删除容器
	containerID := resp.ID
	defer func() {
		err := c.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
		if err != nil {
			panic(err)
		}
	}()

	// 启动容器
	err = c.ContainerStart(ctx, containerID, container.StartOptions{})
	if err != nil {
		panic(err)
	}
	// 查询随机映射选择的端口
	inspRes, err := c.ContainerInspect(ctx, containerID)
	if err != nil {
		panic(err)
	}
	hostPort := inspRes.NetworkSettings.Ports[containerPort][0]
	*mongoURI = fmt.Sprintf("mongodb://%s:%s", hostPort.HostIP, hostPort.HostPort)

	return m.Run()
}
