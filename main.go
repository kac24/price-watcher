package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// PriceResponse is a "Struct". It tells Go how to read the JSON from the API.
type PriceResponse struct {
	Bitcoin struct {
		USD float64 `json:"usd"`
	} `json:"bitcoin"`
}

func checkPrice() {
	// 1. Get the data from CoinGecko
	resp, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd")
	if err != nil {
		fmt.Println("Error fetching price:", err)
		return
	}
	defer resp.Body.Close()

	// 2. Decode the JSON into our Struct
	var data PriceResponse
	json.NewDecoder(resp.Body).Decode(&data)

	price := data.Bitcoin.USD
	fmt.Printf("Current Bitcoin Price: $%.2f\n", price)

	// 3. Logic: If price is low, send a notification via ntfy.sh
	// Change 70000 to whatever price you want to test with
	if price < 100000 { 
		sendNotification(price)
	}
}

func sendNotification(price float64) {
	topic := "my-secret-go-alerts-for-bitcoin-fortesting" // Change this to something unique!
	msg := fmt.Sprintf("Vibe Check: Bitcoin is at $%.2f! Time to look at the charts.", price)
	
	_, err := http.Post("https://ntfy.sh/"+topic, "text/plain", strings.NewReader(msg))
	if err != nil {
		fmt.Println("Error sending notification:", err)
	} else {
		fmt.Println("Notification sent to ntfy.sh/" + topic)
	}
}

func main() {
	fmt.Println("Starting the Vibe Watcher...")
	for {
		checkPrice()
		// Wait 60 seconds before checking again
		time.Sleep(60 * time.Second)
	}
}