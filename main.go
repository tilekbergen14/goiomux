package main

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
)

// type so interface {
//     roomname() string
// }

type room struct {
    name string
}

// func (r room) roomname() string {
//     return r.name
// }



func main() {
	router := mux.NewRouter()
	router.Use(enableCORS)
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		// s.SetContext("")
		fmt.Println("New connection")
		return nil
	})

	server.OnEvent("/", "create", func(s socketio.Conn, roomname string) {
		fmt.Println(roomname)
		s.Join(roomname)
		r := room{name: roomname}
		s.SetContext(r)
		fmt.Println(s.Context())
	})


	server.OnEvent("/", "msg", func(s socketio.Conn, msg string) {
		fmt.Println("New message:", s.Rooms())
		// s.Emit("reply", msg)
		server.BroadcastToRoom("", "room1", "reply", msg)
	})

	

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		// server.Remove(s.ID())
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {

		// Add the Remove session id. Fixed the connection & mem leak
		// server.Remove(s.ID())
		fmt.Println("closed =>", reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()
	router.Handle("/socket.io/", server)
	router.Handle("/", http.FileServer(http.Dir("./asset")))

	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")
	
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	r.Header.Del("Origin")
	
	next.ServeHTTP(w, r)

})}