package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"hw3/chat-service/internal/converter"
	"hw3/chat-service/internal/service"
	"hw3/chat-service/internal/transport/dto"
	"hw3/chat-service/internal/transport/kafka"
)

const (
	lastMessagesCnt = 10
	chatTopic       = "chat"
)

type ChatHandler struct {
	kafkaHandler     *kafka.ChatHandler
	chatService      service.ChatService
	messageConverter converter.MessageConverter
	upgrader         websocket.Upgrader
	validate         *validator.Validate
	clients          map[*websocket.Conn]struct{}
}

func NewChatHandler(
	service service.ChatService,
	upgrader websocket.Upgrader,
	messageConverter converter.MessageConverter,
	kafkaHandler *kafka.ChatHandler,
) *ChatHandler {
	return &ChatHandler{
		kafkaHandler:     kafkaHandler,
		chatService:      service,
		messageConverter: messageConverter,
		upgrader:         upgrader,
		validate:         validator.New(validator.WithRequiredStructEnabled()),
		clients:          make(map[*websocket.Conn]struct{}),
	}
}

func (h *ChatHandler) Chat(w http.ResponseWriter, r *http.Request) {
	connection, _ := h.upgrader.Upgrade(w, r, nil)
	defer connection.Close()

	h.clients[connection] = struct{}{}
	defer delete(h.clients, connection)

	h.sendLastMessages(connection, lastMessagesCnt)

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

func (h *ChatHandler) sendLastMessages(conn *websocket.Conn, count int64) {
	ctx := context.Background()

	lastMessages, err := h.chatService.GetLastMessages(ctx, count)
	if err != nil {
		h.messageToClient(conn, []byte(""))
		log.Fatal(err)
		return
	}
	log.Println(lastMessages)

	responseMessages := h.messageConverter.MapSliceDomainToResponse(lastMessages)
	messagesJson, err := json.Marshal(responseMessages)
	if err != nil {
		h.messageToClient(conn, []byte("Error to send message"))
		log.Println(err)
	} else {
		h.messageToClient(conn, messagesJson)
	}
}

func (h *ChatHandler) readMessage(message []byte) (dto.MessageInfoRequest, error) {
	log.Println(message)
	var info dto.MessageInfoRequest
	err := json.Unmarshal(message, &info)

	if err != nil {
		return dto.MessageInfoRequest{}, err
	}

	err = h.validate.Struct(info)
	if err != nil {
		return dto.MessageInfoRequest{}, err
	}

	return info, nil
}

func (h *ChatHandler) handleMessage(infoRequest dto.MessageInfoRequest) {
	err := h.kafkaHandler.ProduceMessage(chatTopic, infoRequest)
	if err != nil {
		log.Printf("error while send message to kafka: %v", err)
	}

	resp := h.messageConverter.MapRequestToResponse(infoRequest)
	messageBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		return
	}
	h.messageToAllClients(messageBytes)
}
