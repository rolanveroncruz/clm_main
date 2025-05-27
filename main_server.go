package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"ph.certs.com/clm_main/middleware"
	"time"
)

type StringResponse struct {
	Message string `json:"Message"`
}

func init() {}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	portStr := ":" + os.Getenv("PORT")

	mux_router := mux.NewRouter()
	mux_router.Handle("/ping", middleware.LoggingMiddleware(http.HandlerFunc(ping)))
	println("Listening on port " + portStr + "...")
	err = http.ListenAndServe(portStr, mux_router)
	if err != nil {
		panic(err)
	}
}

func ping(w http.ResponseWriter, _ *http.Request) {
	//_, err := w.Write([]byte("Hello, World!"))
	//if err != nil {
	//	return
	//}
	time.Sleep(1 * time.Second)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := StringResponse{
		Message: "Hello, World!",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		return
	}
}
