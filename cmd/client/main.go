package main

import (
	"bufio"
	"fmt"
	"net/http"
	"sync"
)

func makeConnection(i int) {
	url := "http://localhost:8080/events"

	// Create an HTTP client
	client := &http.Client{}

	// Create the GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v", err)
	}

	// Set appropriate headers for SSE
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Bad response status: %s", resp.Status)
	}

	// Create a scanner to read the response line by line
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		// Print each line received from the server
		fmt.Println("Received:", line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading response: %v", err)
	}
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			makeConnection(i)
		}(i)
	}
	wg.Wait()
}
