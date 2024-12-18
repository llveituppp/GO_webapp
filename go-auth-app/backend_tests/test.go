package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	data := map[string]string{
		"email":    "test5@example.com",
		"password": "test",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}

	resp, err := http.Post("http://localhost:8000/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error making POST request:", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	log.Println("Response:", result)
}
