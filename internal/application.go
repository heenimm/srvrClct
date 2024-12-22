package internal

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"serverCalc/pkg"
	"strings"
)

type Config struct {
	AddressOfPort string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.AddressOfPort = os.Getenv("PORT")
	if config.AddressOfPort == "" {
		config.AddressOfPort = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func NewApplication() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

func (a *Application) Run() error {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		text = strings.TrimSpace(text)
		if text == "exit" {
			return nil
		}

		result, err := pkg.Calculate(text)
		if err != nil {
			return err
		}
		log.Panicln(result)
	}
	return nil
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

type CalculationResponse struct {
	Result *float64 `json:"result,omitempty"`
	Error  string   `json:"error,omitempty"`
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	var response CalculationResponse

	if r.Method != http.MethodPost {
		http.Error(w, "method  must be POST 'http://localhost:8082/api/v1/calculate'", http.StatusMethodNotAllowed)
		return
	}

	var request CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error": Expression is not valid}`, http.StatusUnprocessableEntity)
		return
	}

	result, err := pkg.Calculate(request.Expression)

	if err != nil && errors.Is(err, pkg.ErrInternalError) {
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
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
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Panicln("Error encoding response:", err)
	}
}

func (a *Application) RunServer() error {
	// Создаём обработчик для маршрута "/calculate"
	http.HandleFunc("/api/v1/calculate", CalculateHandler)
	if err := http.ListenAndServe(":"+a.config.AddressOfPort, nil); err != nil {
		return pkg.FailedToStartServer
	}
	return nil
}
