package main

import (
	"time"

	"github.com/gorilla/websocket"
)

/*
クライアント・モデル
*/
type client struct {
	socket *websocket.Conn // websocketへのコネクション
	send   chan *message   // メッセージ
	room   *chatroom       // 所属するチャットルーム
}

/*
websocketに書き出されたメッセージを読み込む。
*/
func (c *client) read() {

	// websocketからjson形式でメッセージを読み出し、forwardチャネルに流す。
	// 読み込みは無限ループで実行される。
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			t := time.Now()
			layout := "2006-01-02 15:04:05"
			msg.Time = t.Format(layout)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

/*

 */
func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
