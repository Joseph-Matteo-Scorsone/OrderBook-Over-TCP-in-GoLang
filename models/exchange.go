package models

import (
	"fmt"
	"net"
	"sync"
)

// Message represents a communication packet in the exchange system
type Message struct {
	From    string // Sender identifier
	Payload []byte // Message content
}

// Exchange represents a messaging exchange/broker for client connections
type Exchange struct {
	name      string                // Name of the exchange/channel
	clients   map[net.Conn]struct{} // Map of connected clients (using empty struct as value for minimal memory)
	broadcast chan Message          // Channel for broadcasting messages to all clients
	mu        sync.Mutex            // Mutex for thread-safe operations
}

// NewExchange creates and initializes a new Exchange instance
func NewExchange(name string) *Exchange {
	return &Exchange{
		name:      name,                        // Set exchange name
		clients:   make(map[net.Conn]struct{}), // Initialize clients map
		broadcast: make(chan Message, 1024),    // Initialize buffered broadcast channel with capacity 1024
	}
}

// Start runs the exchange's main message broadcasting loop
func (ex *Exchange) Start() {
	// Continuously process messages from the broadcast channel
	for msg := range ex.broadcast {
		ex.mu.Lock() // Lock for thread-safe access to clients map

		// Iterate through all connected clients
		for client := range ex.clients {
			// Write formatted message to client (sender: payload)
			_, err := client.Write([]byte(fmt.Sprintf("%s: %s\n", msg.From, string(msg.Payload))))
			if err != nil {
				// Handle write errors by logging, closing connection, and removing client
				fmt.Println("write error: ", err)
				client.Close()
				delete(ex.clients, client)
			}
		}
		ex.mu.Unlock() // Unlock after all clients have been processed
	}
}

// Join adds a new client connection to the exchange
func (ex *Exchange) Join(client net.Conn) {
	ex.mu.Lock()         // Lock for thread-safe map modification
	defer ex.mu.Unlock() // Ensure unlock happens when function returns

	// Add client to the map using empty struct as value
	ex.clients[client] = struct{}{}

	// Log client connection and send welcome message
	fmt.Printf("Client %s joined channel %s\n", client.RemoteAddr(), ex.name)
	client.Write([]byte("Welcome to channel " + ex.name + "\n"))
}

// Broadcast queues a message to be sent to all connected clients
func (ex *Exchange) Broadcast(msg Message) {
	ex.broadcast <- msg // Send message to broadcast channel for processing
}
