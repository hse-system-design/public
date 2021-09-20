package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"time"
)


func handleRoot(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello from server"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "plain/text")
}

type HTTPHandler struct {
	storage map[string]string
}

func getRandomKey() string {
	alphabet := []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890")
	rand.Shuffle(len(alphabet), func(i, j int) {
		alphabet[i], alphabet[j] = alphabet[j], alphabet[i]
	})
	id := string(alphabet[:5])
	return id
}

type PutRequestData struct {
	Url string `json:"url"`
}

type PutResponseData struct {
	Key string `json:"key"`
}

func (h *HTTPHandler) handlePutUrl(w http.ResponseWriter, r *http.Request) {
	/*
	{
		"url": "..."
	}
	 */
	var data PutRequestData

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUrlKey := getRandomKey()
	h.storage[newUrlKey] = data.Url
	//  http://my.site.com/bdfhfd

	response := PutResponseData{
		Key: newUrlKey,
	}
	rawResponse, _ := json.Marshal(response)

	_, err = w.Write(rawResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func (h *HTTPHandler) handleGetUrl(w http.ResponseWriter, r *http.Request) {
	// redirect
}

func main() {
	r := mux.NewRouter()

	handler := &HTTPHandler{
		storage: make(map[string]string),
	}

	r.HandleFunc("/", handleRoot).Methods("GET", "POST")
	r.HandleFunc("/{shortUrl:\\w{5}}", handler.handleGetUrl)
	r.HandleFunc("/api/urls", handler.handlePutUrl)


	// add url  ->
	// get url  -> /

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
