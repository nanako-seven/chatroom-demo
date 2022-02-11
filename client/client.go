package main

import (
	"chatroom/api"
	"chatroom/wsutil"
	"context"
	"fmt"
	"time"

	"nhooyr.io/websocket"
)

type Client struct {
	ServerURL string
	Ctx       context.Context
}

const timeout = 100 * time.Second

func (c *Client) Run() {
	ctx, cancel := context.WithTimeout(c.Ctx, timeout)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, c.ServerURL, nil)
	if err != nil {
		return
	}
	defer conn.Close(websocket.StatusInternalError, "the sky is falling")
	defer fmt.Println("与服务器断开连接")

	fmt.Println("连接服务器成功")
	ws := &wsutil.WebsocketWrapper{
		Conn: conn,
		Ctx:  c.Ctx,
	}
	for {
		var username string
		fmt.Printf("%s", "敢问尊姓大名？")
		fmt.Scanln(&username)
		req := &api.ClientUserEnter{
			Username: username,
		}
		err = ws.WriteJSON(req, timeout)
		if err != nil {
			return
		}
		resp := &api.ServerUserEnter{}
		err = ws.ReadJSON(resp, timeout)
		if err != nil {
			return
		}
		if resp.OK {
			fmt.Println("欢迎进入聊天室")
			break
		}
		fmt.Println("名字重复了，换一个吧")
	}

	go func() {
		for {
			var msg string
			fmt.Scanln(&msg)
			req := &api.ClientChatMessage{
				Content: msg,
			}
			err = ws.WriteJSON(req, timeout)
			if err != nil {
				return
			}
		}
	}()

	for {
		resp := &api.ServerChatMessage{}
		err = ws.ReadJSON(resp, timeout)
		if err != nil {
			break
		}
		switch resp.Type {
		case api.Enter:
			fmt.Printf("用户[%s]进入了聊天室\n", resp.Username)
		case api.Leave:
			fmt.Printf("用户[%s]离开了聊天室\n", resp.Username)
		case api.Message:
			fmt.Printf("[%s]: %s\n", resp.Username, resp.Content)
		}
	}
	conn.Close(websocket.StatusNormalClosure, "")
}
