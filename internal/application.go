package internal

import (
	"bufio"
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

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalculateHandler)
	if err := http.ListenAndServe(":"+a.config.AddressOfPort, nil); err != nil {
		return pkg.FailedToStartServer
	}
	return nil
}
