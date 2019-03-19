package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

/* websocket用の変数 */
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

/*
チャットルーム・モデル
*/
type chatroom struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

/*
chatroomをhttp.handleに適合させる。
ここでは以下のことを実装する。
	・websocketの開設
	・clientの生成
*/
func (c *chatroom) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	/* websocketの開設 */
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("websocketの開設に失敗しました。:", err)
	}

	/* クライアントの生成 */
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   c,
	}

	// チャットルームのjoinチャネルにアクセスし、クライアントを入室させる。最後には必ず退室させる。
	c.join <- client
	defer func() {
		c.leave <- client
	}()

	go client.write()
	client.read()
}

/*
チャットルームを生成する
*/
func newRoom() *chatroom {
	t := time.Now()
	layout := "2006-01-02 15:04:05"
	fmt.Println("chatroom が生成されました。:", t.Format(layout))
	return &chatroom{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

/*
チャットルームを起動する
*/
func (c *chatroom) run() {

	// チャットルームは無限ルームで起動させる
	for {
		// チャネルの動きを監視し、処理を決定する
		select {

		/* joinチャネルに動きがあった場合(クライアントの入室) */
		case client := <-c.join:
			// クライアントmapのbool値を真にする
			c.clients[client] = true
			fmt.Printf("クライアントが入室しました。現在 %x 人のクライアントが存在します。\n",
				len(c.clients))

		/* leaveチャネルに動きがあった場合(クライアントの退室) */
		case client := <-c.leave:
			// クライアントmapから対象クライアントを削除する
			delete(c.clients, client)
			fmt.Printf("クライアントが退室しました。現在 %x 人のクライアントが存在します。\n",
				len(c.clients))

		/* forwardチャネルに動きがあった場合(メッセージの受信) */
		case msg := <-c.forward:
			fmt.Printf("メッセージを受信しました。 : %s\n", msg)
			fmt.Printf("全てのクライアントへメッセージを送信します。\n")
			// 存在するクライアント全てに対してメッセージを送信する
			for target := range c.clients {
				select {
				case target.send <- msg:
					fmt.Println("メッセージの送信に成功しました。")
				default:
					fmt.Println("メッセージの送信に失敗しました。")
					delete(c.clients, target)
				}
			}
		}
	}
}
