package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"2sem/hw1/homework/internal/domain"
	"2sem/hw1/homework/internal/service"
	"2sem/hw1/homework/internal/transport/dto"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	service  service.ChatService
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]struct{}
	validate *validator.Validate
}

func NewChatHandler(service service.ChatService, upgrader websocket.Upgrader) *ChatHandler {
	return &ChatHandler{
		service:  service,
		upgrader: upgrader,
		clients:  make(map[*websocket.Conn]struct{}),
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (h *ChatHandler) Chat(w http.ResponseWriter, r *http.Request) {
	connection, _ := h.upgrader.Upgrade(w, r, nil)
	defer connection.Close()

	h.clients[connection] = struct{}{}
	defer delete(h.clients, connection)

	h.sendLastMessages(connection, 10)

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break
		}

		info, err := h.readMessage(message)
		if err != nil {
			log.Println(err)
			continue
		}

		go h.handleMessage(info)
		go printMessage(fmt.Sprintf("%s:%s", info.Nickname, info.Message))
	}
}

func (h *ChatHandler) messageToAllClients(message []byte) {
	for conn := range h.clients {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (h *ChatHandler) messageToClient(conn *websocket.Conn, message []byte) {
	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println(err)
	}
}

func printMessage(message string) {
	fmt.Println(message)
}

func (h *ChatHandler) sendLastMessages(conn *websocket.Conn, count int) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	lastMessages, err := h.service.GetLastMessages(ctx, count)
	if err != nil {
		h.messageToClient(conn, []byte("Error while getting last 10 messages"))
		log.Println(err)
		return
	}
	socketMessages := make([]dto.MessageInfoDTO, 0)

	for _, message := range lastMessages {
		messageDTO := dto.MessageInfoDTO{Nickname: message.Nickname, Message: message.Message, Time: message.Time}
		socketMessages = append(socketMessages, messageDTO)
	}

	messagesJson, err := json.Marshal(socketMessages)
	if err != nil {
		h.messageToClient(conn, []byte("Error to send message"))
		log.Println(err)
	} else {
		h.messageToClient(conn, messagesJson)
	}
}

func (h *ChatHandler) readMessage(message []byte) (dto.MessageInfoDTO, error) {
	var info dto.MessageInfoDTO
	err := json.Unmarshal(message, &info)

	if err != nil {
		return dto.MessageInfoDTO{}, err
	}

	err = h.validate.Struct(info)
	if err != nil {
		return dto.MessageInfoDTO{}, err
	}

	return info, nil
}

func (h *ChatHandler) handleMessage(infoDTO dto.MessageInfoDTO) {
	info, err := h.service.InsertMessage(context.Background(), domain.MessageInfo{
		Nickname: infoDTO.Nickname,
		Message:  infoDTO.Message,
	})
	if err != nil {
		log.Println(err)
		return
	}

	messageDTO := dto.MessageInfoDTO{Nickname: info.Nickname, Message: info.Message, Time: info.Time}
	messageBytes, err := json.Marshal(messageDTO)
	if err != nil {
		log.Println(err)
		return
	}
	h.messageToAllClients(messageBytes)
}
