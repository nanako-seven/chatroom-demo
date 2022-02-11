package wsutil

import (
	"context"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type WebsocketWrapper struct {
	Conn *websocket.Conn
	Ctx  context.Context
}

func (w *WebsocketWrapper) ReadJSON(i interface{}, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(w.Ctx, timeout)
	defer cancel()
	return wsjson.Read(ctx, w.Conn, i)
}

func (w *WebsocketWrapper) WriteJSON(i interface{}, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(w.Ctx, timeout)
	defer cancel()
	return wsjson.Write(ctx, w.Conn, i)
}
