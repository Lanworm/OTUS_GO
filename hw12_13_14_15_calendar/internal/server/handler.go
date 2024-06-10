package server

import "net/http"

type Handler struct{}

func (h *Handler) Hello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("hello handler"))
}
