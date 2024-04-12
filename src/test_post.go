package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var timePost int64

type value struct {
	BannerId int `json:"banner_id"`
}

func sendRequest(banner []byte) string {
	value := value{}
	start := time.Now()
	req, err := http.NewRequest("POST", "http://localhost:8080/banner", bytes.NewBuffer(banner))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return ""
	}
	req.Header.Set("token", "admin_token")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		response, err := io.ReadAll(resp.Body)
		if err == nil {
			json.Unmarshal(response, &value)
		}
	}
	end := time.Now()
	timePost += end.Sub(start).Milliseconds()
	//}
	//chanData <- value.BannerId
	return resp.Status
}

func main() {
	banners := make([]map[string]interface{}, 0, 1000)
	for i := 1; i < 500; i++ {
		for j := 1; j <= 10; j++ {
			count := rand.Intn(4) + 1
			tags := make([]int32, 0, count)
			for k := int32(j); k < int32(j+count); k++ {
				tags = append(tags, k)
			}
			j += count + 1
			requestBody := map[string]interface{}{
				"is_active":  true,
				"feature_id": i,
				"tag_ids":    tags,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
			}
			banners = append(banners, requestBody)
		}
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	total := len(banners)
	success := atomic.Int32{}

	for _, banner := range banners {
		wg.Add(1)
		time.Sleep(1 * time.Millisecond)
		go func(banner map[string]interface{}, success *atomic.Int32) {
			data, _ := json.Marshal(banner)
			time.Sleep(1 * time.Millisecond)
			res := sendRequest(data)
			if res == "201 Created" {
				mu.Lock()
				success.Add(1)
				fmt.Printf("Успешных вставок: %d \n\n", success.Load())
				mu.Unlock()
			} else {
				mu.Lock()
				println(res, banner)
				mu.Unlock()
			}
		}(banner, &success)
		wg.Done()
	}
	wg.Wait()
	fmt.Printf("Выполнено %d вставок \n", total)
	fmt.Printf("Время выполнения %d вставок %vms\n", total, timePost)
	fmt.Printf("Среднее время выполнения одной вставки %vms\n", float64(timePost)/1000)
	fmt.Printf("Успешных вставок: %d \n\n", success.Load())

}
