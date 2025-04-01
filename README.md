#TCP Exchange
TCP Exchange is a Go-based trading exchange system that implements a server-client architecture with order book matching functionality. It supports multiple order types and allows clients to connect via TCP to join exchanges, submit trades, and manage orders.
##Features
Order Book Matching: Supports market orders, Fill-or-Kill (FOK), and Good-Til-Cancelled (GTC) orders with continuous matching.

TCP Server: Handles multiple client connections, allowing them to join/leave exchanges and submit trades.

Message Broadcasting: Broadcasts messages to connected clients within an exchange.

Thread Safety: Uses mutex locks to ensure safe concurrent access to shared resources.

Test Suite: Includes a test function to demonstrate order book functionality.

##Project Structure

TCP-Exchange/
├── models/
│   ├── orderbook.go  # OrderBook struct and matching logic
│   ├── exchange.go   # Exchange struct for message broadcasting
│   └── server.go     # Server struct for TCP connections and client management
└── main.go           # Main entry point and test function

###Key Components
OrderBook (models/orderbook.go):
Manages bids and asks in a trading order book.

Matches orders every 20ms in a background goroutine.

Supports different order types (market, FOK, GTC).

Exchange (models/exchange.go):
Handles client connections within a named exchange.

Broadcasts messages to all connected clients.

Server (models/server.go):
Listens for TCP connections.

Processes client commands (JOIN, LEAVE, TRADE).

Maintains multiple exchanges with their order books.

Main (main.go):
Launches the TCP server on port 8080.

Includes a test function to demonstrate OrderBook functionality.

#Prerequisites
Go 1.16 or higher

#Installation
Clone the repository:
bash

git clone <repository-url>
cd TCP-Exchange

Ensure all dependencies are in place (standard library only):
The project uses only Go standard library packages.

Build and run:
bash

go build
./TCP-Exchange

#Usage
Running the Server:
Start the server by running the compiled binary or go run main.go.

The server listens on :8080 by default.

Connecting as a Client:
Use a TCP client (e.g., telnet or a custom client) to connect to localhost:8080.

Example with telnet:
bash

telnet localhost 8080

Client Commands:
JOIN <exchange_name>: Join an exchange (creates it if it doesn't exist).

LEAVE <exchange_name>: Leave an exchange.

TRADE <exchange_name> <side> <order_type> <price> <size> <ticker>: Submit a trade.
Example: TRADE main buy market 100.00 10 QQQ

Testing OrderBook:
The TestOrderBook() function in main.go runs automatically on startup.

It demonstrates various order types and matching scenarios, printing results to stdout.

##Order Types
Market: Executes immediately at the best available price.

FOK (Fill-or-Kill): Must be filled entirely immediately or cancelled.

GTC (Good-Til-Cancelled): Remains active until matched or server stops.

##Notes
The server runs indefinitely until interrupted (e.g., Ctrl+C).

Error handling is basic; production use would require more robust error management.

The matching engine runs every 20ms; this interval could be adjusted in orderbook.go.

Future Improvements
Add authentication for clients.

Implement order cancellation.

Add persistence for GTC orders across server restarts.

Enhance error handling and logging.

Add support for more order types (e.g., limit orders with partial fills).

##License
This project is licensed under the MIT License - see the LICENSE file for details (Note: Add a LICENSE file if you choose to use this).

