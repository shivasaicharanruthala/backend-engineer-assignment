package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	store "github/shivasaicharanruthala/backend-engineer-assignment/data"
	"github/shivasaicharanruthala/backend-engineer-assignment/handler"
	"github/shivasaicharanruthala/backend-engineer-assignment/log"
	"github/shivasaicharanruthala/backend-engineer-assignment/service"
)

func init() {
	var EnvFilePaths = []string{".env", "/opt/receipts.dev.env", "/opt/receipts.prod.env"}

	for _, envFilePath := range EnvFilePaths {
		err := godotenv.Load(envFilePath)
		if err != nil {
			fmt.Printf("Error loading %v file\n", envFilePath)
		} else {
			break
		}
	}
}

func main() {
	logFilePath := os.Getenv("LOG_FILE_PATH")

	// Initialize Logger
	logger, err := log.NewCustomLogger(logFilePath)
	if err != nil {
		fmt.Printf("Error initiating logger with error %v", err.Error())
		return
	}

	lm := log.Message{Level: "INFO", Msg: "Logger initialized successfully"}
	logger.Log(&lm)

	// Store Layer
	receiptsStore := store.New(logger)

	// Service Layer
	receiptsSvc := service.New(logger, receiptsStore)

	// Handler Layer
	receiptsHandler := handler.New(logger, receiptsSvc)

	// Setup router using mux
	router := mux.NewRouter().StrictSlash(true)
	router.MethodNotAllowedHandler = http.HandlerFunc(MethodNotImplementedHandler)

	// Health check Route
	router.HandleFunc("/v1/health", receiptsHandler.Health).Methods("GET")

	// Receipts Routes
	router.HandleFunc("/v1/receipts/{id}/points", receiptsHandler.Get).Methods("GET")
	router.HandleFunc("/v1/receipts/process", receiptsHandler.Insert).Methods("POST")

	// Start the server
	port := os.Getenv("PORT")
	server := fmt.Sprintf(":%s", port)

	lm = log.Message{Level: "INFO", Msg: fmt.Sprintf("Receipts Server starting to listen on port %v", port)}
	logger.Log(&lm)

	err = http.ListenAndServe(server, router)
	if err != nil {
		lm = log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Initializing receipts server to listen on port %v with error %v", port, err.Error())}
		logger.Log(&lm)
		return
	}
}

func MethodNotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}
