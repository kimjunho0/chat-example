package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// 이 서버에 연결할 클라이언트에 대한 배열
var clients []websocket.Conn

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		//연결 되면 clients 배열에 conn 된 클라이언트를 push 함
		clients = append(clients, *conn)
		//무한 루프로 클라이언트로부터 오는 메시지를 읽는다
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			//서버를 돌리는 콘솔에서도 보내는 메시지가 보이게끔 함
			remoteAddr := conn.RemoteAddr().String()[6:]
			fmt.Printf("%s sent : %s\n", remoteAddr, string(msg))

			messageWithAddr := fmt.Sprintf("%s: %s", remoteAddr, msg)
			//clients 배열에 있는 모든 클라이언트들에게 메시지를 보냄 = 접속한 모든 사용자에게 메시지가 간다.
			for _, client := range clients {
				if err = client.WriteMessage(msgType, []byte(messageWithAddr)); err != nil {
					return
				}
			}
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.ListenAndServe(":8080", nil)
}
