package connection

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func SocketIO() *socketio.Server {
	server := socketio.NewServer(nil)

	//    LOGGING NECESSARY    //

	server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("SocketIO client connected:", s.ID())
		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		log.Println("SocketIO Client disconnected:", s.ID(), "-", msg)
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("SocketIO error:", e)
	})

	server.OnEvent("/", "join-room", func(s socketio.Conn, roomId string) {
		log.Println("SocketIO client", s.ID(), "join-room", roomId)
		server.JoinRoom("/", roomId, s)
	})

	server.OnEvent("/", "leave-room", func(s socketio.Conn, roomId string) {
		log.Println("SocketIO client", s.ID(), "leave-room", roomId)
		server.LeaveRoom("/", roomId, s)
	})

	server.OnEvent("/", "send-message", func(s socketio.Conn, data string) {
		log.Println("SocketIO client", s.ID(), "send-message", data)
		for _, room := range s.Rooms() {
			server.BroadcastToRoom("/", room, "receive-message", data)
			log.Println("SocketIO client", s.ID(), "broadcasting receive-message", data, "to room", room)
		}
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatal("socketio listen error:", err)
		}
	}()

	return server
}
