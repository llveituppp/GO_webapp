package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	// Замените <ваш_токен> на полученный токен
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RpbmdtZW93QGV4YW1wbGUuY29tIiwiZXhwIjoxNzM0Mjc0Mjc4fQ.X92BvcCAbFbUzTnpFiSZCYZjdQySFgixSQOJc2UckGI"

	// Создаём GET-запрос
	req, err := http.NewRequest("GET", "http://localhost:8000/protected/example", nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Добавляем заголовок Authorization
	req.Header.Set("Authorization", "Bearer "+token)

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Читаем и выводим ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response: %s\n", string(body))
}
