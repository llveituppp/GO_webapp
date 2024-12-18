package main

//
import (
	"bytes"
	"encoding/json"
	"io"
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

	resp, err := http.Post("http://localhost:8000/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error making POST request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Server returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	log.Println("Raw response:", string(body))

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal("Error decoding JSON response:", err)
	}

	token := result["token"].(string)

	purchaseData := map[string]interface{}{
		"product_id": 2,
		"quantity":   3,
	}

	purchaseJSON, err := json.Marshal(purchaseData)
	if err != nil {
		log.Fatal("Error encoding purchase data:", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8000/purchase", bytes.NewBuffer(purchaseJSON))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("Error making POST request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Server returned non-200 status: %d", resp.StatusCode)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	log.Println("Raw response:", string(body))

	var purchaseResponse map[string]string
	if err := json.Unmarshal(body, &purchaseResponse); err != nil {
		log.Fatal("Error decoding purchase response:", err)
	}

	log.Println("Purchase response:", purchaseResponse)
}
