package main

import (
	"github.com/gorilla/mux"
	"log"
	"miniurl/handlers"
	"miniurl/storage/mongostorage"
	"net/http"
	"os"
	"time"
)

func NewServer() *http.Server {
	r := mux.NewRouter()

	mongoUrl := os.Getenv("MONGO_URL")
	mongoStorage := mongostorage.NewStorage(mongoUrl)

	handler := &handlers.HTTPHandler{
		Storage: mongoStorage,
	}

	r.HandleFunc("/", handlers.HandleRoot).Methods("GET", "POST")
	r.HandleFunc("/{shortUrl:\\w{5}}", handler.HandleGetUrl).Methods(http.MethodGet)
	r.HandleFunc("/api/urls", handler.HandlePostUrl)

	return &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func main() {
	srv := NewServer()
	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
