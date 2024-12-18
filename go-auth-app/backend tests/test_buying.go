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
		"email":    "test5@example.com", // Ваш тестовый email
		"password": "test",              // Ваш тестовый пароль
	}

	// Преобразуем данные в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}

	// Отправляем запрос на авторизацию
	resp, err := http.Post("http://localhost:8000/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error making POST request:", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Server returned non-200 status: %d", resp.StatusCode)
	}

	// Читаем ответ и извлекаем токен
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

	// Извлекаем токен
	token := result["token"].(string)

	// Данные для покупки товара (например, товар с ID 1 и количество 2)
	purchaseData := map[string]interface{}{
		"product_id": 2, // ID товара
		"quantity":   3, // Количество
	}

	// Преобразуем данные для покупки в JSON
	purchaseJSON, err := json.Marshal(purchaseData)
	if err != nil {
		log.Fatal("Error encoding purchase data:", err)
	}

	// Создаем запрос для покупки товара
	req, err := http.NewRequest("POST", "http://localhost:8000/purchase", bytes.NewBuffer(purchaseJSON))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	// Добавляем JWT токен в заголовок
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос на покупку товара
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("Error making POST request:", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Server returned non-200 status: %d", resp.StatusCode)
	}

	// Читаем ответ
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	log.Println("Raw response:", string(body))

	// Декодируем JSON-ответ (сообщение о результате покупки)
	var purchaseResponse map[string]string
	if err := json.Unmarshal(body, &purchaseResponse); err != nil {
		log.Fatal("Error decoding purchase response:", err)
	}

	// Выводим результат покупки
	log.Println("Purchase response:", purchaseResponse)
}
