package internal

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"serverCalc/handlers"
	"serverCalc/middleware"
	"serverCalc/pkg"
	"strconv"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
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
	db     *sqlx.DB
}

func NewApplication(db *sqlx.DB) *Application {
	return &Application{
		config: ConfigFromEnv(),
		db:     db,
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
	var wg sync.WaitGroup

	computingPowerStr := os.Getenv("COMPUTING_POWER")
	computingPower, err := strconv.Atoi(computingPowerStr)
	if err != nil || computingPower <= 0 {
		computingPower = 1
	}

	go pkg.GenerateTask()
	mux := http.NewServeMux()
	mux.Handle("/api/v1/calculate", middleware.JWTAuthMiddleware(http.HandlerFunc(AddExpressionHandler)))
	mux.Handle("/api/v1/expressions", middleware.JWTAuthMiddleware(http.HandlerFunc(GetExpressionsHandler)))
	mux.Handle("/api/v1/expressions/", middleware.JWTAuthMiddleware(http.HandlerFunc(GetExpressionByIDHandler)))

	mux.HandleFunc("/api/v1/login", handlers.LoginHandler(a.db))
	mux.HandleFunc("/api/v1/register", handlers.RegisterHandler(a.db))
	mux.HandleFunc("/internal/task", InternalTaskHandler)

	if err := http.ListenAndServe(":"+a.config.AddressOfPort, mux); err != nil {
		return pkg.FailedToStartServer
	}

	for i := 0; i < computingPower; i++ {
		wg.Add(1)
		go worker(&wg, "http://localhost:"+a.config.AddressOfPort)
	}

	go func() {
		wg.Wait()
	}()

	return nil
}
