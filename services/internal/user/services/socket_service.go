package services

import (
	"log"
	"services/internal/user/dto"

	"github.com/gofiber/contrib/websocket"
)

type SocketService struct {
}

func NewSocketService() SocketService {
	return SocketService{}
}

func (s *SocketService) ReadMessage(conn *websocket.Conn) (code int, message []byte) {
	mt, msg, err := conn.ReadMessage()
	if err != nil {
		s.WriteMessage(conn, websocket.TextMessage, []byte(err.Error()))
	}
	return mt, msg
}
func (s *SocketService) ParseMessage(conn *websocket.Conn, msg []byte) *dto.NewMessageDTO {
	unmarshalledData, err := dto.NewWebSocketReadMessageDTO(msg)
	if err != nil {
		s.WriteMessage(conn, websocket.TextMessage, []byte(err.Error()))
	}

	return unmarshalledData
}

func (s *SocketService) WriteMessage(conn *websocket.Conn, code int, message []byte) {

	if err := conn.WriteMessage(code, message); err != nil {
		log.Println("Socket Write Error", err)
	}
}
