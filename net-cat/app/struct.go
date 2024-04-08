package netcat

import (
	"net"
	"sync"
)

type Client struct {
	Name     string
	Messages chan string
}

type ChatHistory struct {
	Messages []string
	MaxSize  int
	sync.Mutex
}

var (
	linux = `
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\\dS"qML
 |    '.       |  \\' Zq
_)      \\.___.,|     .'
\\____   )MMMMMP|   .'
     '-'       '--'`

	clients   = make(map[net.Addr]Client)
	clientsMx sync.Mutex
	history   = ChatHistory{
		Messages: make([]string, 0),
		MaxSize:  50, // Максимальное количество сохраняемых сообщений истории чата
	}
	semaphore = make(chan struct{}, 10) // Ограничитель на 10 клиентов
)
