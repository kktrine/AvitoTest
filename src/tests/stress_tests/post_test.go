package stress_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var timePost int64

//type value struct {
//	BannerId int `json:"banner_id"`
//}

func sendPostRequest(banner []byte) string {
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
	if err == nil && resp != nil {
		defer resp.Body.Close()
	} else {
		return "ERROR"
	}
	end := time.Now()
	timePost += end.Sub(start).Milliseconds()
	return resp.Status
}

func TestAdd(t *testing.T) {
	n := 1000
	banners := make([]map[string]interface{}, 0, n)
	tags := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 1; i < n+1; i++ {
		requestBody := map[string]interface{}{
			"is_active":  true,
			"feature_id": i,
			"tag_ids":    tags,
			"content": map[string]string{
				"title": "some_title111",
				"text":  "some_text",
				"url":   "some_url",
			},
		}
		if rand.Intn(9) == 1 {
			requestBody["is_active"] = false
		}
		banners = append(banners, requestBody)
	}
	banners[998]["is_active"] = true
	banners[999]["is_active"] = false
	var wg sync.WaitGroup
	var mu sync.Mutex
	total := len(banners)
	success := atomic.Int32{}
	timePost = 0
	for _, banner := range banners {
		wg.Add(1)
		time.Sleep(1 * time.Millisecond)
		go func(banner map[string]interface{}, success *atomic.Int32) {
			defer wg.Done()
			data, _ := json.Marshal(banner)
			time.Sleep(1 * time.Millisecond)
			res := sendPostRequest(data)
			if res == "201 Created" {
				mu.Lock()
				success.Add(1)
				mu.Unlock()
			} else {
				mu.Lock()
				fmt.Println(res, banner)
				fmt.Println(banner)
				mu.Unlock()
			}
		}(banner, &success)

	}
	wg.Wait()
	fmt.Printf("Выполнено %d запросов на вставку \n", total)
	fmt.Printf("Время выполнения %d запросов %vms\n", total, timePost)
	fmt.Printf("Среднее время выполнения одного запроса %vms\n", float64(timePost)/1000)
	numSuccess := success.Load()
	fmt.Printf("Успешных запросов: %d = %v %% \n\n", numSuccess, float64(numSuccess)*100./float64(total))

}
