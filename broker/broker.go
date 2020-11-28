package broker

import (
	"log"
	"net/http"
	"strings"
)

type Broker struct {

	// Create a map of clients, the keys of the map are the channels
	// over which we can push messages to attached clients. (The values
	// are just booleans and are meaningless.)
	//
	Clients map[chan string]bool

	// Channel into which new clients can be pushed
	NewClients chan chan string

	// Channel into which disconnected clients should be pushed
	DefunctClients chan chan string

	// Channel into which messages are pushed to be broadcast out
	// to attahed clients.
	Messages chan []byte
}

func (b *Broker) Start() {
	go func() {
		for {
			select {
			case s := <-b.NewClients:
				b.Clients[s] = true
			case s := <-b.DefunctClients:
				delete(b.Clients, s)
				close(s)
			case msg := <-b.Messages:
				for s := range b.Clients {
					s <- string(msg)
				}
			}
		}
	}()
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Make sure that the writer supports flushing.
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan string)

	b.NewClients <- messageChan

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		// Remove this client from the map of attached clients
		// when `EventHandler` exits.
		b.DefunctClients <- messageChan
		log.Println("HTTP connection just closed.")
	}()

	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	for {
		// Read from our messageChan.
		msg, open := <-messageChan
		if !open {
			break
		}
		w.Write(formatSSE("message", string(msg)))
		f.Flush()
	}
}

func formatSSE(event string, data string) []byte {
	eventPayload := "event: " + event + "\n"
	dataLines := strings.Split(data, "\n")
	for _, line := range dataLines {
		eventPayload = eventPayload + "data: " + line + "\n"
	}
	return []byte(eventPayload + "\n")
}
