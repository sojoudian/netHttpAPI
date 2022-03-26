package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Vote struct {
	Name     string `json:"name"`
	Vote     string `json:"vote"`
	Hostname string `json:"hostname"`
	ID       string `json:"id"`
}

type voteAppHandler struct {
	sync.Mutex
	store map[string]Vote
}

func (h *voteAppHandler) VotottingApp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
	default:
		h.get(w, r)
		return
	}
}

func (h *voteAppHandler) get(w http.ResponseWriter, r *http.Request) {
	response := make([]Vote, len(h.store))

	h.Lock()
	i := 0
	for _, vote := range h.store {
		response[i] = vote
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *voteAppHandler) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	var vote Vote

	err = json.Unmarshal(bodyBytes, &vote)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("Content-Type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need 'cotent-type': 'application/json', but got '%s'", ct)))
		return
	}

	vote.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[vote.ID] = vote
	defer h.Unlock()

}

func newVoteAppHandler() *voteAppHandler {
	return &voteAppHandler{
		store: map[string]Vote{
			// "id_!23": Vote{
			// 	Name:     "maziar",
			// 	Vote:     "Cats",
			// 	Hostname: "xyz123",
			// 	ID:       "id_123",
			// },
		},
	}
}

func main() {
	voteAppHandler := newVoteAppHandler()
	http.HandleFunc("/", voteAppHandler.VotottingApp)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}

}
