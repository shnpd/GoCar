package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// websocket建立在http协议之上，所以需要先建立一个http连接
	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 升级到websocket
	u := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	c, err := u.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("upgrade error: %v\n", err)
		return
	}
	defer c.Close()

	done := make(chan struct{})
	go func() {
		for {
			m := make(map[string]interface{})
			err := c.ReadJSON(&m)
			if err != nil {
				fmt.Printf("read error: %v\n", err)
				// 如果不是正常错误则打印错误信息
				if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					fmt.Printf("unexpected read error: %v\n", err)
				}
				// 通知主goroutine关闭
				done <- struct{}{}
				break
			}
			fmt.Printf("recv: %v\n", m)
		}
	}()

	i := 0
	for {
		// 每200ms发送一条消息，如果在200ms内收到done信号则退出
		select {
		case <-time.After(200 * time.Millisecond):
		case <-done:
			return
		}

		i++
		err := c.WriteJSON(map[string]string{
			"message": "hello",
			"msgid":   strconv.Itoa(i),
		})
		if err != nil {
			fmt.Printf("write error: %v\n", err)
		}
	}
}
