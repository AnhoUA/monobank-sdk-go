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

	fmt.Println("\nMCC Info Example (English by default):")
	if info, ok := client.MCC.Get("0742"); ok {
		fmt.Printf("MCC: %s\nGroup: %s\nDescription: %s\n", info.MCC, info.GroupName, info.ShortDescription)
	}

	fmt.Println("\nMCC Info Example (Ukrainian):")
	if info, ok := client.MCC.GetUk("0742"); ok {
		fmt.Printf("MCC: %s\nGroup: %s\nDescription: %s\n", info.MCC, info.GroupName, info.ShortDescription)
	}
}
