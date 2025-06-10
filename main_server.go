package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"ph.certs.com/clm_main/auth"
	"ph.certs.com/clm_main/certs"
	"ph.certs.com/clm_main/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	portStr := ":" + os.Getenv("PORT")

	muxRouter := mux.NewRouter()
	muxRouter.Handle("/ping", middleware.LoggingMiddleware(http.HandlerFunc(ping)))
	muxRouter.Handle("/postping", middleware.LoggingMiddleware(http.HandlerFunc(postPing)))
	muxRouter.Handle("/api/get_server_cert", middleware.LoggingMiddleware(middleware.JWTMiddleware(http.HandlerFunc(certs.GetServerCert))))
	muxRouter.Handle("/api/login", middleware.LoggingMiddleware(middleware.CorsMiddleware(http.HandlerFunc(auth.Login))))
	println("Listening on port " + portStr + "...")
	err = http.ListenAndServe(portStr, muxRouter)
	if err != nil {
		panic(err)
	}
}

func ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	type StringResponse struct {
		Message string `json:"Message"`
	}

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

func postPing(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	type RequestData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	type StringResponse struct {
		Message string `json:"Message"`
	}
	var data RequestData
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error Decoding JSON", http.StatusBadRequest)
		return
	}
	name := data.Name
	email := data.Email

	response := StringResponse{
		Message: "Hello, " + name + " with email " + email,
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
