package capoeira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type HTTPTransport struct {
	endpoints []string
	// channel keys are of the form "from->to"
	receivedMessages map[string]chan interface{}
	server           *http.Server
	port             int
	lock             sync.RWMutex
}

func NewHTTPTransport(endpoints []string) *HTTPTransport {
	t := &HTTPTransport{
		endpoints:        endpoints,
		receivedMessages: make(map[string]chan interface{}),
		port:             8080,
	}
	for _, endpoint := range endpoints {
		if _, ok := t.receivedMessages[endpoint]; !ok {
			t.receivedMessages[endpoint] = make(chan interface{}, 1) // buffered channel to avoid deadlock
		}
	}
	t.StartServer()
	return t
}

func (t *HTTPTransport) Send(from, to string, data any) {
	fmt.Println("HTTPTransport sending from", from, "to", to, "data:", data)
	payload := map[string]any{
		"from": from,
		"to":   to,
		"data": data,
	}
	fmt.Println("Payload:", payload)

	b, err := json.Marshal(payload)

	if err != nil {
		fmt.Printf("Error marshaling payload: %v\n", err)
		return
	}
	resp, err := http.Post("http://localhost:8080/message", "application/json", bytes.NewBuffer(b))
	if err != nil {
		fmt.Printf("Error sending HTTP request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Non-OK HTTP status: %s\n", resp.Status)
	}
}

func (t *HTTPTransport) Receive(from, at string) interface{} {
	t.lock.RLock()
	defer t.lock.RUnlock()
	fmt.Printf("Receiving on %s...\n", at)
	val := <-t.receivedMessages[at]
	fmt.Printf("Received at %s from %s: %v of type %T\n", at, from, val, val)
	return val
}

func (t *HTTPTransport) Locations() []string {
	return t.endpoints
}

// StartServer starts an HTTP server to listen for incoming messages on the given port
func (t *HTTPTransport) StartServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusBadRequest)
			return
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}
		fmt.Println("Received payload:", payload)
		to, _ := payload["to"].(string)
		// put the received message onto the channel for this pair of from/to locations
		t.lock.RLock()
		t.receivedMessages[to] <- payload["data"]
		fmt.Printf("Wrote %v to channel %v\n", payload["data"], to)
		t.lock.RUnlock()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	addr := fmt.Sprintf(":%d", t.port)
	t.server = &http.Server{Addr: addr, Handler: mux}
	go func() {
		if err := t.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()
	fmt.Printf("HTTPTransport server started on port %d\n", t.port)
	return nil
}

// StopServer gracefully stops the HTTP server
func (t *HTTPTransport) StopServer() error {
	if t.server != nil {
		return t.server.Close()
	}
	return nil
}
