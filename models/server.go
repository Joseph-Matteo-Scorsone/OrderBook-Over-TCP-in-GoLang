package models

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Server represents a trading server managing client connections and order books
type Server struct {
	listenAddr string                // Network address to listen on
	ln         net.Listener          // Network listener for incoming connections
	quitch     chan struct{}         // Channel to signal server shutdown
	clients    map[net.Conn]string   // Map of client connections to their exchange names
	exchanges  map[string]*OrderBook // Map of exchange names to their order books
	mu         sync.Mutex            // Mutex for thread-safe operations
}

// NewServer creates and initializes a new Server instance
func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,                  // Set listening address
		quitch:     make(chan struct{}),         // Initialize quit channel
		clients:    make(map[net.Conn]string),   // Initialize clients map
		exchanges:  make(map[string]*OrderBook), // Initialize exchanges map
	}
}

// StartServer begins the server's operation
func (s *Server) StartServer() error {
	ln, err := net.Listen("tcp", s.listenAddr) // Start TCP listener
	if err != nil {
		return err // Return error if listening fails
	}
	defer ln.Close() // Ensure listener closes when function returns

	s.ln = ln // Store listener in server struct

	go s.acceptLoop() // Start accepting connections in separate goroutine

	<-s.quitch // Wait for quit signal

	return nil // Return nil on successful shutdown
}

// acceptLoop continuously accepts new client connections
func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept() // Accept incoming connection
		if err != nil {
			fmt.Println("accept error:", err) // Log acceptance errors
			continue
		}

		go s.handleClient(conn) // Handle each client in a separate goroutine
	}
}

// handleClient manages communication with a connected client
func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close() // Ensure connection closes when function returns

	fmt.Println("new connection to server:", conn.RemoteAddr()) // Log new connection

	reader := bufio.NewReader(conn) // Create buffered reader for client input
	for {
		msg, err := reader.ReadString('\n') // Read message until newline
		if err != nil {
			fmt.Println("read error:", err) // Log read errors
			return
		}

		msg = strings.TrimSpace(msg) // Remove leading/trailing whitespace
		// Handle JOIN command
		if strings.HasPrefix(msg, "JOIN ") {
			exchangeName := strings.TrimPrefix(msg, "JOIN ")
			s.joinExchange(exchangeName, conn)

			// Handle LEAVE command
		} else if strings.HasPrefix(msg, "LEAVE ") {
			exchangeName := strings.TrimPrefix(msg, "LEAVE ")
			s.leaveExchange(exchangeName, conn)

			// Handle TRADE command
		} else if strings.HasPrefix(msg, "TRADE ") {
			parts := strings.SplitN(msg, " ", 7) // Split into max 7 parts
			if len(parts) != 7 {
				conn.Write([]byte("Invalid message format\n")) // Validate message format
				continue
			}

			exchangeName := parts[1]
			order := orderFromParts(parts) // Parse order from message parts

			s.trade(exchangeName, order, conn) // Process trade

		} else {
			conn.Write([]byte("Unknown command\n")) // Handle unknown commands
			continue
		}
	}
}

// orderFromParts constructs an Order from message parts
func orderFromParts(parts []string) *Order {
	side := parts[2]                             // Buy or sell
	orderType := parts[3]                        // Type of order
	price, _ := strconv.ParseFloat(parts[4], 64) // Parse price (ignoring error)
	size, _ := strconv.Atoi(parts[5])            // Parse size (ignoring error)
	ticker := parts[6]                           // Trading pair/symbol

	return &Order{
		Price:     price,              // Set order price
		Ticker:    ticker,             // Set trading ticker
		Size:      size,               // Set order size
		OrderType: GetType(orderType), // Convert string to order type
		Side:      side,               // Set buy/sell side
		Timestamp: time.Now(),         // Set current timestamp
	}
}

// joinExchange adds a client to an exchange
func (s *Server) joinExchange(exchangeName string, conn net.Conn) {
	s.mu.Lock()         // Lock for thread safety
	defer s.mu.Unlock() // Ensure unlock on return

	// Create new order book if exchange doesn't exist
	if _, exists := s.exchanges[exchangeName]; !exists {
		s.exchanges[exchangeName] = NewOrderBook()
	}
	s.clients[conn] = exchangeName                                         // Associate client with exchange
	conn.Write([]byte(fmt.Sprintf("Joined exchange: %s\n", exchangeName))) // Confirm join
}

// leaveExchange removes a client from an exchange
func (s *Server) leaveExchange(exchangeName string, conn net.Conn) {
	s.mu.Lock()         // Lock for thread safety
	defer s.mu.Unlock() // Ensure unlock on return

	// Note: Current implementation doesn't check if client was in the exchange
	if _, exists := s.exchanges[exchangeName]; !exists {
		conn.Write([]byte(fmt.Sprintf("Left exchange: %s\n", exchangeName)))
	}
	delete(s.clients, conn)                                              // Remove client from clients map
	conn.Write([]byte(fmt.Sprintf("Left exchange: %s\n", exchangeName))) // Confirm leave
}

// trade processes a trading order in an exchange
func (s *Server) trade(exchangeName string, order *Order, conn net.Conn) {
	s.mu.Lock()         // Lock for thread safety
	defer s.mu.Unlock() // Ensure unlock on return

	exchange, exists := s.exchanges[exchangeName] // Get exchange
	if !exists {
		conn.Write([]byte(fmt.Sprintf("Not in exchange %s\n", exchangeName))) // Check existence
		return
	}

	exchange.AddOrder(*order)                                                     // Add order to exchange's order book
	conn.Write([]byte(fmt.Sprintf("Order added to %s exchange\n", exchangeName))) // Confirm trade
}
