package main

import (
	"context"
	"fmt"
	"time"
)

type paramKey struct{}

func main() {
	// 使用paramKey而不是"param"字符串,作为强类型在使用时不需要手动拼写字符串
	c := context.WithValue(context.Background(), paramKey{}, "abc")
	// context一旦创建后就无法改变，这里context.WithTimeout返回一个新的context再将其重新赋值给c
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	go mainTask(c)

	var cmd string
	for {
		fmt.Scan(&cmd)
		if cmd == "c" {
			cancel()
		}
	}
}

func mainTask(c context.Context) {
	fmt.Printf("main task started with param%q\n", c.Value(paramKey{}))
	go func() {
		c1, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		smallTask(c1, "task1", 9*time.Second)
	}()
	smallTask(c, "task2", 8*time.Second)
}

// 约定context始终作为第一个参数
func smallTask(c context.Context, name string, d time.Duration) {
	fmt.Printf("%s started with param %q\n", name, c.Value(paramKey{}))
	select {
	// 任务需要d才能完成
	case <-time.After(d):
		fmt.Printf("%s done\n", name)
	case <-c.Done():
		// 如果5s的时间到了，或者调用了cancel函数，c.Done()就会接收到信号
		fmt.Printf("%s cancelled\n", name)
	}
}
