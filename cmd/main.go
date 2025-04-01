package main

import (
	"TCP-Exchange/models" // Import the models package containing server and order book logic
	"fmt"
	"time"
)

// main is the entry point of the program
func main() {
	TestOrderBook() // Run the order book test function

	// Start the TCP server
	server := models.NewServer(":8080")          // Create a new server listening on port 8080
	if err := server.StartServer(); err != nil { // Start the server and handle any errors
		fmt.Println("Error starting server:", err)
	}
}

// TestOrderBook tests the functionality of the OrderBook implementation
func TestOrderBook() {
	ob := models.NewOrderBook() // Create a new order book instance
	TICKER := "QQQ"             // Define the trading ticker symbol

	// Test 1: Matching market orders at same price
	ob.AddOrder(models.Order{
		Price:     100.00,
		Ticker:    TICKER,
		Size:      10,
		OrderType: models.GetType("market"), // Market order type
		Side:      "buy",
	})
	ob.AddOrder(models.Order{
		Price:     100.00,
		Ticker:    TICKER,
		Size:      10,
		OrderType: models.GetType("market"),
		Side:      "sell",
	})

	time.Sleep(30 * time.Millisecond) // Wait for matching to occur (matching runs every 20ms)

	// Print current state of order book
	fmt.Printf("Bids: %v\n", ob.Bids)
	fmt.Printf("Asks: %v\n", ob.Asks)
	fmt.Printf("\n")

	// Test 2: Non-matching market orders at different prices
	ob.AddOrder(models.Order{
		Price:     99.00,
		Ticker:    TICKER,
		Size:      10,
		OrderType: models.GetType("market"),
		Side:      "buy",
	})
	ob.AddOrder(models.Order{
		Price:     101.00,
		Ticker:    TICKER,
		Size:      10,
		OrderType: models.GetType("market"),
		Side:      "sell",
	})

	time.Sleep(30 * time.Millisecond) // Wait for matching attempt

	fmt.Printf("Bids: %v\n", ob.Bids)
	fmt.Printf("Asks: %v\n", ob.Asks)
	fmt.Printf("\n")

	// Test 3: Matching Fill-or-Kill (FOK) orders with equal sizes
	ob.AddOrder(models.Order{
		Price:     100.00,
		Ticker:    TICKER,
		Size:      10,
		OrderType: models.GetType("fok"), // Fill-or-Kill order type
		Side:      "buy",
	})
	ob.AddOrder(models.Order{
		Price:     100.00,
		Ticker:    TICKER,
		Size:      10,
		OrderType: models.GetType("fok"),
		Side:      "sell",
	})

	time.Sleep(30 * time.Millisecond) // Wait for matching

	fmt.Printf("Bids: %v\n", ob.Bids)
	fmt.Printf("Asks: %v\n", ob.Asks)
	fmt.Printf("\n")

	// Test 4: Non-matching FOK orders with different sizes
	ob.AddOrder(models.Order{
		Price:     100.00,
		Ticker:    TICKER,
		Size:      5, // Smaller size
		OrderType: models.GetType("fok"),
		Side:      "buy",
	})
	ob.AddOrder(models.Order{
		Price:     100.00,
		Ticker:    TICKER,
		Size:      10, // Larger size
		OrderType: models.GetType("fok"),
		Side:      "sell",
	})

	time.Sleep(30 * time.Millisecond) // Wait for matching attempt

	fmt.Printf("Bids: %v\n", ob.Bids)
	fmt.Printf("Asks: %v\n", ob.Asks)
	fmt.Printf("\n")

	// Test 5: Good-Til-Cancelled (GTC) orders that don't match
	ob.AddOrder(models.Order{
		Price:     98.00,
		Ticker:    TICKER,
		Size:      90,
		OrderType: models.GetType("gtc"), // Good-Til-Cancelled order type
		Side:      "buy",
	})
	ob.AddOrder(models.Order{
		Price:     102.00,
		Ticker:    TICKER,
		Size:      30,
		OrderType: models.GetType("gtc"),
		Side:      "sell",
	})

	time.Sleep(30 * time.Millisecond) // Wait for matching attempt

	fmt.Printf("Bids: %v\n", ob.Bids)
	fmt.Printf("Asks: %v\n", ob.Asks)
	fmt.Printf("\n")

	// Test 6: Stop matching and check GTC order persistence
	ob.StopMatching() // Stop the order matching process

	time.Sleep(30 * time.Millisecond) // Wait to ensure matching has stopped

	fmt.Printf("Bids: %v\n", ob.Bids) // Should show remaining GTC orders
	fmt.Printf("Asks: %v\n", ob.Asks)
	fmt.Printf("\n")
}
