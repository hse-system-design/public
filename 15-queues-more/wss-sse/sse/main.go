package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/sse", sseHandler)
	http.HandleFunc("/", rootHandler)

	panic(http.ListenAndServe(":8080", nil))
}

//go:embed index.html
var indexHTML string

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "%s", indexHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for idx := 0; true; idx ++ {
		message := fmt.Sprintf("Message number %v", idx)
		_, err := fmt.Fprintf(w, "data: %v\n\n", message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		fmt.Printf("Send msg `%s`\n", message)
		time.Sleep(500 * time.Millisecond)
	}
}

