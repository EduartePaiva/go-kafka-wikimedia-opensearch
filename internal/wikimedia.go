package internal

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	WIKIMEDIA_URL = "https://stream.wikimedia.org/v2/stream/recentchange"
)

func TestWikimedia() {

	// Create a new HTTP request
	req, err := http.NewRequest("GET", WIKIMEDIA_URL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	// Set a custom User-Agent header
	req.Header.Set("User-Agent", "wikimedia-pet-project/1.0")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("-------------------------------------------------")
		fmt.Fprintln(os.Stdout, line)
		fmt.Println("-------------------------------------------------")
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading stream: %v", err)
	}
}
