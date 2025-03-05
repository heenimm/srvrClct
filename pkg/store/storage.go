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
	mutex       = sync.Mutex{}
	path        = "expressions.json"
)

func init() {
	if err := loadFromFile(); err != nil {
		log.Println("Ошибка загрузки данных из файла:", err)
	}
}

func AddExpression(id, exp string) error {
	mutex.Lock()
	defer mutex.Unlock()

	expressions[id] = pkg.ExpressionResponse{
		ID:     id,
		Status: "success",
		Result: exp,
	}
	log.Println("success при сохранении JSON:")

	return saveToFile()
}

func GetExpressions() ([]pkg.ExpressionResponse, error) {
	if err := loadFromFile(); err != nil {
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

func saveToFile() error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	log.Println("Ошибка при кодировании JSON:", err)
	return encoder.Encode(expressions)
}

func loadFromFile() error {
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
