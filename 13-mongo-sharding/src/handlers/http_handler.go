package handlers

import (
	"encoding/json"
	"fmt"
	"miniurl/ratelimit"
	"miniurl/storage"
	"net/http"
	"strings"
	"time"
)

func HandleRoot(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello from server"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "plain/text")
}

func NewHTTPHandler(
	storage storage.Storage,
	limiterFactory *ratelimit.Factory,
	indexMaintainers []storage.IndexMaintainer,
) *HTTPHandler {
	return &HTTPHandler{
		Storage:          storage,
		indexMaintainers: indexMaintainers,
		postLimit:        limiterFactory.NewLimiter("post_url", 10*time.Second, 2),
		getLimit:         limiterFactory.NewLimiter("get_url", 1*time.Minute, 10),
	}
}

type HTTPHandler struct {
	Storage storage.Storage

	indexMaintainers []storage.IndexMaintainer

	postLimit *ratelimit.Limiter
	getLimit  *ratelimit.Limiter
}

type PutRequestData struct {
	Url string `json:"url"`
}

type PutResponseData struct {
	Key string `json:"key"`
}

func (h *HTTPHandler) HandlePostUrl(rw http.ResponseWriter, r *http.Request) {
	//canDo, err := h.postLimit.CanDoAt(r.Context(), time.Now())
	//if err != nil {
	//	http.Error(rw, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//if !canDo {
	//	http.Error(rw, "rate limit exceeded", http.StatusTooManyRequests)
	//	return
	//}

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
	//canDo, err := h.getLimit.CanDoAt(r.Context(), time.Now())
	//if err != nil {
	//	http.Error(rw, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//if !canDo {
	//	http.Error(rw, "rate limit exceeded", http.StatusTooManyRequests)
	//	return
	//}

	key := strings.Trim(r.URL.Path, "/")

	url, err := h.Storage.GetURL(r.Context(), storage.URLKey(key))
	if err != nil {
		http.NotFound(rw, r)
		return
	}
	http.Redirect(rw, r, string(url), http.StatusPermanentRedirect)
}

func (h *HTTPHandler) CreateIndices(rw http.ResponseWriter, r *http.Request) {
	for _, maintainer := range h.indexMaintainers {
		if err := maintainer.EnsureIndices(r.Context()); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	_, _ = rw.Write([]byte("All indices are successfully created"))
}
