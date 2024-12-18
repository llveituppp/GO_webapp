package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func main() {
	// Данные для входа
	data := map[string]string{
		"email":    "test5@example.com",
		"password": "test",
	}

	// Преобразуем данные в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}

	// Отправляем запрос
	resp, err := http.Post("http://localhost:8000/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error making POST request:", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Server returned non-200 status: %d", resp.StatusCode)
	}

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	log.Println("Raw response:", string(body))

	// Декодируем JSON-ответ
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal("Error decoding JSON response:", err)
	}

	log.Println("Response:", result)
}
