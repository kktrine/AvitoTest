package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	banners := make([]map[string]interface{}, 0, 1000)
	tags := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 1; i < 1000; i++ {
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
		//if rand.Intn(9) == 1 {
		//	requestBody["is_active"] = false
		//}
		banners = append(banners, requestBody)
	}
	//for i := 1; i < 1000; i++ {
	//	for j := 1; j <= 10; j++ {
	//		//count := rand.Intn(4) + 1
	//		count := rand.Intn(4) + 1
	//		tags := make([]int32, 0, count)
	//		for k := int32(j); k < int32(j+count); k++ {
	//			tags = append(tags, k)
	//		}
	//		j += count + 1
	//		requestBody := map[string]interface{}{
	//			"is_active":  true,
	//			"feature_id": i,
	//			"tag_ids":    tags,
	//			"content": map[string]string{
	//				"title": "some_title",
	//				"text":  "some_text",
	//				"url":   "some_url",
	//			},
	//		}
	//		banners = append(banners, requestBody)
	//	}
	//}
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
