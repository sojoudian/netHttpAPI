package main

import "net/http"

type Vote struct {
	Name     string `json:"name"`
	Vote     string `json:"vote"`
	Hostname string `json:"hostname"`
	ID       string `json:"id"`
}

type voteAppHandler struct {
	store map[string]Vote
}

func (h *voteAppHandler) get(w http.ResponseWriter, r *http.Request) {

}

func newVoteAppHandler() *voteAppHandler {
	return &voteAppHandler{
		store: map[string]Vote{},
	}
}

func main() {
	voteAppHandler := newVoteAppHandler()
	http.HandleFunc("/", voteAppHandler.get)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}

}
