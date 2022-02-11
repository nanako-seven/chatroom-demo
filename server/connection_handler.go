package main

import (
	"chatroom/api"
	"chatroom/wsutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"nhooyr.io/websocket"
)

type ConnectionHandler struct {
	Server *ChatServer
}

const timeout = 100 * time.Second

func (h *ConnectionHandler) HandleConnection(c *gin.Context) {
	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}
	defer func() {
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			conn.Close(websocket.StatusNormalClosure, "")
		} else {
			conn.Close(websocket.StatusInternalError, "")
		}
	}()

	ws := &wsutil.WebsocketWrapper{
		Conn: conn,
		Ctx:  c.Request.Context(),
	}

	var username string
	var ch chan ChatServerMessage
	req := &api.ClientUserEnter{}
	resp := &api.ServerUserEnter{
		OK: false,
	}
	for {
		err := ws.ReadJSON(req, timeout)
		if err != nil {
			return
		}
		ch, err = h.Server.UserEnter(req.Username)
		if err == nil {
			username = req.Username
			resp.OK = true
			err = ws.WriteJSON(resp, timeout)
			if err != nil {
				return
			}
			defer h.Server.UserLeave(req.Username)
			break
		}
		err = ws.WriteJSON(resp, timeout)
		if err != nil {
			return
		}
	}

	go writeClient(ws, ch)

	for {
		req := &api.ClientChatMessage{}
		err = ws.ReadJSON(req, timeout)
		if err != nil {
			return
		}
		h.Server.Broadcast(username, req.Content)
	}
}

func writeClient(w *wsutil.WebsocketWrapper, ch chan ChatServerMessage) {
	for {
		msg := <-ch
		resq := &api.ServerChatMessage{
			Username: msg.Username,
			Content:  msg.Content,
		}
		switch msg.Type {
		case ChatMessageTypeEnter:
			resq.Type = api.Enter
		case ChatMessageTypeLeave:
			resq.Type = api.Leave
		case ChatMessageTypeMessage:
			resq.Type = api.Message
		}
		err := w.WriteJSON(resq, timeout)
		if err != nil {
			return
		}
	}
}
