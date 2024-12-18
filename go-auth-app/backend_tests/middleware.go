package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// test
func main() {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RpbmdtZW93QGV4YW1wbGUuY29tIiwiZXhwIjoxNzM0Mjc0Mjc4fQ.X92BvcCAbFbUzTnpFiSZCYZjdQySFgixSQOJc2UckGI"

	req, err := http.NewRequest("GET", "http://localhost:8000/protected/example", nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response: %s\n", string(body))
}
