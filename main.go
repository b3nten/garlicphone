package main

import (
	"6enten/garlicphone/messages"
	"fmt"
	"net/http"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

var fs = http.FileServer(http.Dir("./web"))

var (
	upgrader = newUpgrader()
)

func newUpgrader() *websocket.Upgrader {
	u := websocket.NewUpgrader()
	u.OnOpen(func(c *websocket.Conn) {
		// echo
		fmt.Println("OnOpen:", c.RemoteAddr().String())
	})
	u.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		// echo
		fmt.Println("OnMessage:", messageType, string(data))
		c.WriteMessage(messageType, data)
	})
	u.OnClose(func(c *websocket.Conn, err error) {
		fmt.Println("OnClose:", c.RemoteAddr().String(), err)
	})
	return u
}

func onWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Upgraded:", conn.RemoteAddr().String())
}

func sendBinary(mux *http.ServeMux) {
	player := &messages.Player{
		Id:   messages.Ptr(uint32(12345)),
		Name: messages.Ptr("PlayerOne"),
		Inventory: &[]messages.Item{
			{ Name: messages.Ptr("Hammer") },
			{ Name: messages.Ptr("Sword") },
		},
	}
	for i := 0; i < 1000; i++ {
		*player.Inventory = append(*player.Inventory, messages.Item{ Name: messages.Ptr(fmt.Sprintf("Item%d", i)) })
	}
	bytes, err := messages.MarshalBinary(player)

	if err != nil {
		panic(err)
	}
	// fmt.Println("Serialized player to binary", bytes)

	mux.HandleFunc("/binary", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(bytes)
	})
}

func main() {
	mux := &http.ServeMux{}
	mux.HandleFunc("/ws", onWebsocket)
	mux.HandleFunc("/", fs.ServeHTTP)
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.HandlerFunc(fs.ServeHTTP)).ServeHTTP(w, r)
	})
	sendBinary(mux)
	server := http.Server{
		Addr:    "localhost:8000",
		Handler: mux,
	}
	fmt.Println("server started at http://localhost:8000")
	fmt.Println("server exit:", server.ListenAndServe())
}
