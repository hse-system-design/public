package handlers

import (
	"encoding/json"
	"fmt"
	"miniurl/storage"
	"net/http"
	"strings"
	"sync"
)

func HandleRoot(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello from server"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "plain/text")
}

type HTTPHandler struct {
	StorageMu sync.RWMutex
	Storage   storage.Storage

}

type PutRequestData struct {
	Url string `json:"url"`
}

type PutResponseData struct {
	Key string `json:"key"`
}

func (h *HTTPHandler) HandlePostUrl(rw http.ResponseWriter, r *http.Request) {
	var data PutRequestData

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}


	newUrlKey, err := h.Storage.PutURL(r.Context(), storage.ShortedURL(data.Url))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	//newUrlKey := generator.GetRandomKey()
	//h.StorageMu.Lock()
	//h.Storage[newUrlKey] = data.Url
	//h.StorageMu.Unlock()
	//  http://my.site.com/bdfhfd

	response := PutResponseData{
		Key: string(newUrlKey),
	}
	rawResponse, _ := json.Marshal(response)

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(rawResponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

}

func (h *HTTPHandler) HandleGetUrl(rw http.ResponseWriter, r *http.Request) {
	key := strings.Trim(r.URL.Path, "/")
	//h.StorageMu.RLock()
	//url, found := h.Storage[key]
	//h.StorageMu.RUnlock()

	url, err := h.Storage.GetURL(r.Context(), storage.URLKey(key))
	if err != nil {
		http.NotFound(rw, r)
		return
	}
	http.Redirect(rw, r, string(url), http.StatusPermanentRedirect)
}

