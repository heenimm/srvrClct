package internal

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"net/http"
	"serverCalc/pkg"
	"serverCalc/pkg/store"
	"strings"
	"sync"
)

var (
	expressionStore = make(map[string]string)
	taskQueue       []pkg.CalculationRequest
	mutex           sync.Mutex
)

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	var response pkg.CalculationResponse

	if r.Method != http.MethodPost {
		response.Error = "метод должен быть POST"
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request pkg.CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error = "Expression is not valid"
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	result, err := pkg.Calculate(request.Expression)
	if err != nil && errors.Is(err, pkg.ErrInternalError) {
		response.Error = "Internal server error"
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			panic(err)
		}
		return
	}

	if err != nil {
		response.Error = "Expression is not valid"
		w.WriteHeader(http.StatusUnprocessableEntity)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			panic(err)
		}
		return
	}

	response.Result = &result
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "метод должен быть POST", http.StatusMethodNotAllowed)
		return
	}

	var request pkg.CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.Expression == "" {
		http.Error(w, "неверный запрос", http.StatusUnprocessableEntity)
		return
	}

	id := uuid.New().String()

	mutex.Lock()
	expressionStore[id] = request.Expression
	store.AddExpression(id, request.Expression)
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pkg.CalculationResponse{ID: id})
}

func GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "метод должен быть GET", http.StatusMethodNotAllowed)
		return
	}

	expressions, err := store.GetExpressions()
	if err != nil {
		http.Error(w, "неверный запрос", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"expressions": expressions,
	})
}

func GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "что-то пошло не так, метод должен быть GET", http.StatusMethodNotAllowed)
		return
	}

	vars := strings.Split(r.URL.Path, "/")
	if len(vars) < 5 || vars[4] == "" {
		http.Error(w, "с таким ID нет записей", http.StatusBadRequest)
		return
	}
	id := vars[4][1:]
	expression, err := store.GetExpressionByID(id)
	if err != nil {
		http.Error(w, "нет такого выражения", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expression)
}

func InternalTaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		postTaskResultHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "метод должен быть GET", http.StatusMethodNotAllowed)
		return
	}

	if len(taskQueue) == 0 {
		http.Error(w, "нет задачи", http.StatusNotFound)
		return
	}

	task, err := store.GetTasks()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"task": task,
	})
}

func postTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "метод должен быть POST", http.StatusMethodNotAllowed)
		return
	}

	var result pkg.TaskResponse
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "невалидные данные", http.StatusUnprocessableEntity)
		log.Print(err)
		return
	}

	if _, exists := expressionStore[result.ID]; !exists {
		http.Error(w, "нет такой задачи", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
