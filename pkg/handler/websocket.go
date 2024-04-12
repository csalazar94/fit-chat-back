package handler

import (
	"io"
	"log"
	"net/http"

	"nhooyr.io/websocket"
)

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"echo", "other"},
	})
	if err != nil {
		log.Printf("%v", err)
		return
	}
	log.Printf("Conexión establecida")
	defer conn.CloseNow()

	log.Printf("Subprotocolos: %v", conn.Subprotocol())

	for {
		messageType, reader, err := conn.Reader(r.Context())
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			log.Printf("Conexión cerrada")
			return
		}
		if err != nil {
			log.Printf("%v", err)
			return
		}

		msg, err := io.ReadAll(reader)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		log.Printf("Tipo de mensaje: %v, Mensaje: %v", messageType, string(msg))
	}
}
