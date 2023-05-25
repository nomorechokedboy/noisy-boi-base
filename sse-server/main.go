package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type SSEConn struct {
	mu      sync.Mutex
	clients map[string][]chan string
}

func NewSSEConn() *SSEConn {
	return &SSEConn{clients: make(map[string][]chan string)}
}

func (p *SSEConn) addClient(id string) *chan string {
	p.mu.Lock()
	defer func() {
		fmt.Println("Clients in add: ", p.clients)
		for k, v := range p.clients {
			fmt.Printf("Key: %s, value: %d\n", k, len(v))
			fmt.Println("Channels from id=", id, v)
		}
		p.mu.Unlock()
	}()

	c, ok := p.clients[id]
	if !ok {
		client := []chan string{make(chan string)}
		p.clients[id] = client
		return &client[0]
	}

	newCh := make(chan string)
	p.clients[id] = append(c, newCh)
	return &newCh
}

func (p *SSEConn) removeClient(id string, conn chan string) {
	p.mu.Lock()
	defer func() {
		fmt.Println("Clients in remove: ", p.clients)
		for k, v := range p.clients {
			fmt.Printf("Key: %s, value: %d", k, len(v))
		}
		p.mu.Unlock()
	}()

	c, ok := p.clients[id]
	if !ok {
		return
	}

	pos := -1

	for i, ch := range c {
		if ch == conn {
			pos = i
		}
	}

	if pos == -1 {
		return
	}

	close(c[pos])
	c = append(c[:pos], c[pos+1:]...)
	if pos == 0 {
		delete(p.clients, id)
	}
}

func (p *SSEConn) broadcast(id string, data, event string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	c, ok := p.clients[id]
	if !ok {
		return
	}

	for _, ch := range c {
		ch <- fmt.Sprintf("event: %s\ndata: %s\n\n", event, data)
	}
}

var sseConn = NewSSEConn()

func getTime(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/time/")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	msg := time.Now().Format("15:04:05")
	sseConn.broadcast(id, msg, "timeEvent")
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/sse/")
	ch := sseConn.addClient(id)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	defer sseConn.removeClient(id, *ch)

	flusher, ok := w.(http.Flusher)
	if !ok {
		fmt.Println("Could not init http.Flusher")
	}

	for {
		select {
		case message := <-*ch:
			fmt.Println("case message... sending message")
			fmt.Println(message)
			fmt.Fprintf(w, message)
			flusher.Flush()
		case <-r.Context().Done():
			fmt.Println("Client closed connection")
			return
		}
	}
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/sse/", sseHandler)
	router.HandleFunc("/time/", getTime)

	log.Fatal(http.ListenAndServe(":3500", router))
}
