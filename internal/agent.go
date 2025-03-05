package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"serverCalc/pkg"
	"time"
)

func worker(orchestratorURL string) {
	//defer wg.Done()

	for {
		resp, err := http.Get(orchestratorURL + "/internal/task")
		if err != nil {
			fmt.Println("Ошибка запроса задачи:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Нет задач...")
			time.Sleep(5 * time.Second)
			continue
		}

		var task pkg.CalculationRequest
		err = json.NewDecoder(resp.Body).Decode(&task)
		resp.Body.Close()
		if err != nil {
			fmt.Println("Ошибка декодирования задачи:", err)
			continue
		}

		result, err := pkg.Calculate(task.Expression)
		if err != nil {
			fmt.Println("Ошибка вычисления:", err)
			continue
		}

		resultPayload := pkg.TaskResponse{ID: task.ID, Result: result}
		body, _ := json.Marshal(resultPayload)

		resp, err = http.Post(orchestratorURL+"/internal/task", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Ошибка отправки результата:", err)
			continue
		}
		resp.Body.Close()
	}
}
