package stress_tests

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var timeGet int64

func sendGetUserBannerRequest(args string) string {
	start := time.Now()
	req, err := http.NewRequest("GET", "http://localhost:8080/user_banner"+args, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return ""
	}
	req.Header.Set("token", "user_token")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}
	end := time.Now()
	//pretty.Print(resp.Body)
	timeGet += end.Sub(start).Milliseconds()
	//}
	//chanData <- value.BannerId
	return resp.Status
}

func TestGetUserBanner(t *testing.T) {
	args := make([]string, 0, 1000)
	for i := 1; i < 1000; i++ {
		str := fmt.Sprintf("?tag_id=%d&feature_id=%d&use_last_revision=%d", rand.Intn(9)+1, rand.Intn(998)+1, rand.Intn(1))
		args = append(args, str)
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	total := len(args)
	success := atomic.Int32{}
	timeGet = 0
	for _, arg := range args {
		wg.Add(1)
		time.Sleep(1 * time.Millisecond)
		go func(arg string, success *atomic.Int32) {
			defer wg.Done()
			time.Sleep(1 * time.Millisecond)
			res := sendGetUserBannerRequest(arg)
			if res == "200 OK" {
				mu.Lock()
				success.Add(1)
				mu.Unlock()
			} else {
				mu.Lock()
				println(res, arg)
				mu.Unlock()
			}
		}(arg, &success)
	}
	wg.Wait()
	fmt.Printf("Выполнено %d запросов \n", total)
	fmt.Printf("Время выполнения %d запросов %vms\n", total, timeGet)
	fmt.Printf("Среднее время выполнения одного запроса %vms\n", float64(timeGet)/1000)
	numSuccess := success.Load()
	fmt.Printf("Успешных вставок: %d = %v %% \n\n", numSuccess, float64(numSuccess)*100./float64(total))

}
