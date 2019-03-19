package main

import (
	"log"
	"net/http"
)

func main() {

	/* ルートへのアクセスに対してハンドラを貼り、chat.htmlをサーブする */
	http.Handle("/", &templateHandler{filename: "/chat.html"})

	/* チャットルームを作成する */
	chatroom := newRoom()

	/* チャットルームへのハンドラを貼る。
	   /room へは、chat.html から遷移する。chatroomに実装されたServeHTTPは
	   websocketの実装とclientの生成を行う。 */
	http.Handle("/room", chatroom)

	/* チャットルームを起動する */
	go chatroom.run()

	/* webサーバを開始する */
	log.Println("webサーバを開始します。")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("webサーバの起動に失敗しました。:", err)
	}
}
