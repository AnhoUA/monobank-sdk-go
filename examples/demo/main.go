package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AnhoUA/monobank-sdk-go/sdk"
)

func main() {
	// Create unified client
	client := sdk.NewClient("") // No token needed for public currency rates

	// Get currency rates via client.Public
	currencies, err := client.Public.GetCurrency(context.Background())
	if err != nil {
		log.Fatalf("Error getting rates: %v", err)
	}

	fmt.Println("Monobank currency rates:")
	for _, c := range currencies {
		if c.CurrencyCodeB == 980 { // Only relative to UAH
			fmt.Printf("Currency A: %d, Sell: %.2f, Buy: %.2f\n", c.CurrencyCodeA, c.RateSell, c.RateBuy)
		}
	}
}
