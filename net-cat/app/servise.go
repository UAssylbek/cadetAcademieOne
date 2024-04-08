package netcat

import (
	"fmt"
	"time"
)

func NotifyNewClient(newClientName string) {
	clientsMx.Lock()
	defer clientsMx.Unlock()

	for _, client := range clients {
		if client.Name != newClientName {
			client.Messages <- fmt.Sprintf("\n%s has joined our chat...\n", newClientName)
			client.Messages <- fmt.Sprintf("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + client.Name + "]:")
		}
	}
}

func BroadcastMessage(msg string, senderName string) {
	clientsMx.Lock()
	defer clientsMx.Unlock()

	for _, client := range clients {
		if client.Name != senderName { // Исключаем отправку сообщения обратно отправителю
			client.Messages <- msg
		}
		Word := "[" + time.Now().Format("2006-01-02 15:04:05") + "][" + client.Name + "]:"
		client.Messages <- Word
	}
}
