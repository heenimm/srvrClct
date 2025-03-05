package store

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"serverCalc/pkg"
	"sync"
)

var (
	expressions = make(map[string]pkg.ExpressionResponse)
	tasks       = make(map[string]pkg.TaskResponse)
	mutex       = sync.Mutex{}
)

func AddExpression(id, exp string) error {
	mutex.Lock()
	defer mutex.Unlock()

	expressions[id] = pkg.ExpressionResponse{
		ID:          id,
		Status:      "success",
		Expressions: exp,
	}
	log.Println("Задача успешно сохранена")
	return saveExpToFile()
}

func GetExpressions() ([]pkg.ExpressionResponse, error) {
	if err := loadFromFile("expressions.json"); err != nil {
		log.Println(err)
		return nil, err
	}

	mutex.Lock()
	defer mutex.Unlock()
	expressionsList := make([]pkg.ExpressionResponse, 0, len(expressions))
	for _, expression := range expressions {
		if expression.Status == "success" {
			expressionsList = append(expressionsList, expression)
		}
	}

	return expressionsList, nil
}

func GetExpressionByID(id string) (pkg.ExpressionResponse, error) {
	mutex.Lock()
	defer mutex.Unlock()

	exp, ok := expressions[id]
	if ok {
		return exp, nil
	}
	return pkg.ExpressionResponse{}, errors.New("expression not found")
}

func AddTask(id, arg1, arg2, operation string) error {
	mutex.Lock()
	defer mutex.Unlock()

	tasks[id] = pkg.TaskResponse{
		ID: id,
		//Arg1:      arg1,
		//Arg2:      arg2,
		//Operation: operation,
	}

	return saveTasksToFile()
}

func GetTasks() ([]pkg.TaskResponse, error) {
	if err := loadTasksFromFile("tasks.json"); err != nil {
		log.Println(err)
		return nil, err
	}

	mutex.Lock()
	defer mutex.Unlock()
	tasksList := make([]pkg.TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		tasksList = append(tasksList, task)
	}

	return tasksList, nil
}

func GetLastTask() (pkg.TaskResponse, error) {
	if len(tasks) == 0 {
		return pkg.TaskResponse{}, errors.New("нет сохраненных задач")
	}

	var lastTask pkg.TaskResponse
	for _, task := range tasks {
		lastTask = task
	}

	return lastTask, nil
}

func saveTasksToFile() error {
	file, err := os.Create("tasks.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	log.Println("Ошибка при кодировании JSON:", err)
	return encoder.Encode(tasks)
}

func loadTasksFromFile(path string) error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&tasks)
}

func saveExpToFile() error {
	file, err := os.Create("expressions.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	log.Println("Ошибка при кодировании JSON:", err)
	return encoder.Encode(expressions)
}

func loadFromFile(path string) error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&expressions)
}
