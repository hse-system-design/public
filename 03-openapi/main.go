package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello from server"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "plain/text")
}

type HTTPHandler struct {
	storageMu sync.RWMutex
	storage   map[string]string
}

var alphabet = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890")

func getRandomKey() string {
	idBytes := make([]byte, 5)
	for i := 0; i < len(idBytes); i++ {
		idBytes[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(idBytes)
}

type PutRequestData struct {
	Url string `json:"url"`
}

type PutResponseData struct {
	Key string `json:"key"`
}

func (h *HTTPHandler) handlePostUrl(rw http.ResponseWriter, r *http.Request) {
	var data PutRequestData

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	newUrlKey := getRandomKey()
	h.storageMu.Lock()
	h.storage[newUrlKey] = data.Url
	h.storageMu.Unlock()
	//  http://my.site.com/bdfhfd

	response := PutResponseData{
		Key: newUrlKey,
	}
	rawResponse, _ := json.Marshal(response)

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(rawResponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

}

func (h *HTTPHandler) handleGetUrl(rw http.ResponseWriter, r *http.Request) {
	key := strings.Trim(r.URL.Path, "/")
	h.storageMu.RLock()
	url, found := h.storage[key]
	h.storageMu.RUnlock()
	if !found {
		http.NotFound(rw, r)
		return
	}
	http.Redirect(rw, r, url, http.StatusPermanentRedirect)
}

func NewServer() *http.Server {
	r := mux.NewRouter()

	handler := &HTTPHandler{
		storage: make(map[string]string),
	}

	r.HandleFunc("/", handleRoot).Methods("GET", "POST")
	r.HandleFunc("/{shortUrl:\\w{5}}", handler.handleGetUrl).Methods(http.MethodGet)
	r.HandleFunc("/api/urls", handler.handlePostUrl)

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
