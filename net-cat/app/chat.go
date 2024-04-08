package netcat

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	semaphore <- struct{}{} // Захватываем слот

	// Приветствие и запрос имени
	conn.Write([]byte("Welcome to TCP-Chat!\n"))
	conn.Write([]byte(linux))
	conn.Write([]byte("\n[ENTER YOUR NAME]: "))
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading name:", err)
		return
	}
	name = strings.TrimSpace(name)
	for len(name) == 0 {
		conn.Write([]byte("Error: empty name try again \n"))
		conn.Write([]byte("[ENTER YOUR NAME]: "))
		name, err = bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading name:", err)
			return
		}
		name = strings.TrimSpace(name)
	}
	// Отправка истории чата новому клиенту
	history.Lock()
	for _, msg := range history.Messages {
		conn.Write([]byte(msg + "\n"))
	}
	history.Unlock()
	conn.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name)))

	// Создание нового клиента
	client := Client{
		Name:     name,
		Messages: make(chan string),
	}

	// Добавление клиента в список
	clientsMx.Lock()
	clients[conn.RemoteAddr()] = client
	clientsMx.Unlock()

	// Отправка сообщений клиенту
	go func() {
		for {
			msg, ok := <-client.Messages
			if !ok {
				return
			}
			conn.Write([]byte(msg))
		}
	}()

	// Уведомляем остальных клиентов о подключении нового клиента
	NotifyNewClient(name)
	// Чтение сообщений от клиента
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if msg != "" {
			BroadcastMessage(fmt.Sprintf("\n"+"[%s][%s]: %s\n", time.Now().Format("2006-01-02 15:04:05"), name, msg), name)
			// Добавление сообщений в историю чата
			history.Lock()
			history.Messages = append(history.Messages, fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), name, msg))
			if len(history.Messages) > history.MaxSize {
				history.Messages = history.Messages[len(history.Messages)-history.MaxSize:]
			}
			history.Unlock()
		} else {
			conn.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name)))
		}
	}

	// Удаляем клиента из списка при отключении
	clientsMx.Lock()
	delete(clients, conn.RemoteAddr())
	clientsMx.Unlock()

	// Освобождаем слот
	<-semaphore

	// Сообщаем остальным клиентам об отключении клиента
	BroadcastMessage(fmt.Sprintf("\n%s has left our chat...\n", name), name)
}
