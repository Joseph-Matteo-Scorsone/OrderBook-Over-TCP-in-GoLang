package models

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// OrderBook represents a trading order book with bids and asks
type OrderBook struct {
	Bids        map[float64][]Order // Map of bid prices to their orders
	Asks        map[float64][]Order // Map of ask prices to their orders
	mu          sync.Mutex          // Mutex for thread-safe operations
	nextOrderID int                 // Counter for generating unique order IDs
	stopChan    chan struct{}       // Channel to signal stopping of matching goroutine
}

// NewOrderBook creates and initializes a new OrderBook instance
func NewOrderBook() *OrderBook {
	ob := &OrderBook{
		Bids:     make(map[float64][]Order), // Initialize bids map
		Asks:     make(map[float64][]Order), // Initialize asks map
		stopChan: make(chan struct{}),       // Initialize stop channel
	}

	// Start the order matching process in a separate goroutine
	go ob.startMatching()

	return ob
}

// startMatching runs a continuous order matching process
func (ob *OrderBook) startMatching() {
	ticker := time.NewTicker(20 * time.Millisecond) // Create ticker for 20ms intervals
	go func() {
		for {
			select {
			case <-ticker.C: // Every 20ms, try to match orders
				ob.MatchOrders()
			case <-ob.stopChan: // When stop signal received, clean up and exit
				ticker.Stop()
				return
			}
		}
	}()
}

// MatchOrders attempts to match bid and ask orders
func (ob *OrderBook) MatchOrders() {
	ob.mu.Lock()         // Lock the order book for thread safety
	defer ob.mu.Unlock() // Ensure unlock happens when function returns

	// Exit if either bids or asks are empty
	if len(ob.Bids) == 0 || len(ob.Asks) == 0 {
		return
	}

	matches := false // Track if any matches occurred
	// Iterate through all bid prices
	for bidPrice, bids := range ob.Bids {
		// Iterate through all ask prices
		for askPrice, asks := range ob.Asks {
			// Check if bid price meets or exceeds ask price
			if bidPrice >= askPrice {
				for i := 0; i < len(bids); i++ {
					bid := bids[i]

					for j := 0; j < len(asks); j++ {
						ask := asks[i] // Note: Should be asks[j], potential bug in original

						// Process limit orders (OrderType == 1)
						if bid.OrderType == 1 || ask.OrderType == 1 {
							vol := min(bid.Size, ask.Size) // Find minimum tradeable volume

							// Handle Fill-or-Kill (FOK) orders
							if bid.FOK || ask.FOK {
								// If volume doesn't match full order size, cancel FOK order
								if vol != bid.Size && vol != ask.Size {
									if bid.FOK {
										bids = append(bids[:i], bids[i+1:]...)
										i--
									} else {
										asks = append(asks[:j], asks[j+1:]...)
										j--
									}
									continue
								}
							}

							if vol > 0 { // If there's a match
								matches = true

								// Reduce sizes of both orders
								bid.Size -= vol
								ask.Size -= vol

								// Log the trade
								fmt.Printf("Fill at price: %.2f, size %d, with order id %d and %d\n",
									askPrice, vol, bid.Id, ask.Id)

								// Clean up ask if fully filled
								if ask.Size == 0 {
									ob.Asks[askPrice] = append(ob.Asks[askPrice][:j], ob.Asks[askPrice][j+1:]...)
									j--
									if len(ob.Asks[askPrice]) == 0 {
										delete(ob.Asks, askPrice)
									}
								}

								// Clean up bid if fully filled
								if bid.Size == 0 {
									bids = append(bids[:i], bids[i+1:]...)
									ob.Bids[bidPrice] = bids
									i--
									if len(ob.Bids[bidPrice]) == 0 {
										delete(ob.Bids, bidPrice)
									}
									break
								}
							}
						}
					}
				}
			}
		}
		if !matches {
			fmt.Printf("No matches found this iteration\n")
		}
	}
}

// StopMatching gracefully stops the order matching process
func (ob *OrderBook) StopMatching() {
	close(ob.stopChan) // Signal the matching goroutine to stop

	ob.mu.Lock()         // Lock for thread safety
	defer ob.mu.Unlock() // Ensure unlock on return

	// Process Good-Til-Cancelled (GTC) orders in Bids
	for price, bids := range ob.Bids {
		newBids := []Order{}
		for _, bid := range bids {
			if bid.GTC { // Keep only GTC orders
				newBids = append(newBids, bid)
			}
		}
		if len(newBids) == 0 {
			delete(ob.Bids, price) // Remove price level if no GTC orders remain
		} else {
			ob.Bids[price] = newBids // Update with only GTC orders
		}
	}

	// Process GTC orders in Asks
	for price, asks := range ob.Asks {
		newAsks := []Order{}
		for _, ask := range asks {
			if ask.GTC { // Keep only GTC orders
				newAsks = append(newAsks, ask)
			}
		}
		if len(newAsks) == 0 {
			delete(ob.Asks, price) // Remove price level if no GTC orders remain
		} else {
			ob.Asks[price] = newAsks // Update with only GTC orders
		}
	}
}

// AddOrder adds a new order to the order book
func (ob *OrderBook) AddOrder(order Order) {
	ob.mu.Lock()         // Lock for thread safety
	defer ob.mu.Unlock() // Ensure unlock on return

	// Assign order ID and timestamp
	order.Id = ob.nextOrderID
	order.Timestamp = time.Now()
	ob.nextOrderID++

	// Set order flags based on type
	if order.OrderType == 3 {
		order.GTC = true // Good-Til-Cancelled
	} else if order.OrderType == 4 {
		order.FOK = true // Fill-or-Kill
	}

	// Add order to appropriate side and sort by timestamp
	if order.Side == "buy" {
		ob.Bids[order.Price] = append(ob.Bids[order.Price], order)
		sort.Slice(ob.Bids[order.Price], func(i, j int) bool {
			return ob.Bids[order.Price][i].Timestamp.Before(ob.Bids[order.Price][j].Timestamp)
		})
	} else if order.Side == "sell" {
		ob.Asks[order.Price] = append(ob.Asks[order.Price], order)
		sort.Slice(ob.Asks[order.Price], func(i, j int) bool {
			return ob.Asks[order.Price][i].Timestamp.Before(ob.Asks[order.Price][j].Timestamp)
		})
	}
}
