package models

import "time"

// Order represents a financial order in the system.
type Order struct {
	Price     float64   // The price at which the order is placed
	Id        int       // Unique identifier for the order
	Ticker    string    // The stock or asset symbol associated with the order
	Size      int       // The number of shares or contracts in the order
	OrderType int       // Numeric representation of the order type (e.g., market, limit, etc.)
	Side      string    // Buy or sell indicator
	Timestamp time.Time // Time when the order was created
	GTC       bool      // Good-Till-Canceled flag, true if the order remains active until executed or canceled
	FOK       bool      // Fill-Or-Kill flag, true if the order must be executed immediately in full or canceled
}

// GetType converts a string representation of an order type into its corresponding integer value.
func GetType(orderType string) int {
	switch orderType {
	case "market":
		return 1 // Market order
	case "limit":
		return 2 // Limit order
	case "gtc":
		return 3 // Good-Till-Canceled order
	case "fok":
		return 4 // Fill-Or-Kill order
	}

	return 0 // Unknown order type
}
