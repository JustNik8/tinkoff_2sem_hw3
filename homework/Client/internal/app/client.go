package app

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	"github.com/gorilla/websocket"
	"hw3/client/config"
	"hw3/client/internal/transport/dto"
)

const (
	clientConfigPath = "client_config.yaml"
)

func RunClient() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.ParseClientConfig(clientConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)
	nickname := readNickname(reader)

	host := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	u := url.URL{Scheme: "ws", Host: host, Path: "/chat"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go writeMessages(c, reader, nickname)
	go readMessages(c)

	<-ctx.Done()
	log.Println("Client shutdown gracefully")
}

func readNickname(reader *bufio.Reader) string {
	fmt.Print("Enter your nickname: ")

	nickname, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)

	}

	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		log.Fatal("You input empty nickname")
	}

	return nickname
}

func writeMessages(c *websocket.Conn, reader *bufio.Reader, nickname string) {
	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		messageDTO := dto.MessageInfoRequest{Nickname: nickname, Message: message}
		messageBytes, err := json.Marshal(messageDTO)
		if err != nil {
			log.Println(err)
			break
		}

		err = c.WriteMessage(websocket.TextMessage, messageBytes)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func readMessages(c *websocket.Conn) {
	err := readLastMessages(c)
	if err != nil {
		log.Fatal(err)
	}

	for {
		messageType, messageBytes, err := c.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			log.Fatal(err)
		}
		var messageDTO dto.MessageInfoResponse

		err = json.Unmarshal(messageBytes, &messageDTO)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s:%s\n", messageDTO.Nickname, messageDTO.Message)
	}
}

func readLastMessages(c *websocket.Conn) error {
	messageType, messageBytes, err := c.ReadMessage()
	if err != nil || messageType == websocket.CloseMessage {
		return err
	}

	messagesDTO := make([]dto.MessageInfoResponse, 0)
	err = json.Unmarshal(messageBytes, &messagesDTO)
	if err != nil {
		return err
	}

	sort.Slice(messagesDTO, func(i, j int) bool {
		return messagesDTO[i].Time.Before(messagesDTO[j].Time)
	})
	for _, message := range messagesDTO {
		fmt.Printf("%s:%s\n", message.Nickname, message.Message)
	}

	return nil
}
