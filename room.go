package chatapp

import (
	"fmt"
	"io"
	"sync"
)

type room struct {
	name    string
	msgch   chan string
	quit    chan struct{}
	clients map[chan<- string]struct{}
	*sync.RWMutex
}

func CreateRoom(name string) *room {
	r := room{
		name:    name,
		msgch:   make(chan string),
		clients: make(map[chan<- string]struct{}),
		RWMutex: new(sync.RWMutex),
		quit:    make(chan struct{}),
	}
	r.Run()
	return &r
}

func (r *room) Run() {
	fmt.Println("starting chat room")
	go func() {
		for msg := range r.msgch {
			r.broadcastMsg(msg)
		}
	}()

}

func (r *room) CLCounts() int {
	return len(r.clients)
}

func (r *room) Addclient(c io.ReadWriteCloser) {
	r.Lock()
	wc, done := StartClient(c, r.msgch, r.quit)
	r.clients[wc] = struct{}{}
	r.Unlock()
	go func() {
		<-done
		r.RemoveClient(wc)
	}()

}

func (r *room) RemoveClient(wc chan<- string) {
	fmt.Println("remove client")
	r.Lock()
	close(wc)
	delete(r.clients, wc)
	r.Unlock()
	select {
	case <-r.quit:
		if len(r.clients) == 0 {
			close(r.msgch)
		}
	default:
	}

}

func (r *room) broadcastMsg(msg string) {
	r.RLock()
	defer r.RUnlock()
	fmt.Println("received message", msg)
	for wc, _ := range r.clients {
		go func(wc chan<- string) {
			fmt.Println("1")
			wc <- msg
		}(wc)
	}
}
