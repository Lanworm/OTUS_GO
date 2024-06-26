package httphandler

import (
	"fmt"
	"github.com/Lanworm/OTUS_GO/final_project/internal/logger"
	"net/http"
	"strings"
)

type Handler struct {
	logger *logger.Logger
}

func NewHandler(
	logger *logger.Logger,

) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) ResizeHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	s := r.URL.String()
	v := strings.Split(s, "/")
	width := v[2]
	height := v[3]
	url := strings.Join(v[4:], "/")
	fmt.Println(width, height)
	fmt.Println(url)
	w.Write([]byte("test"))
}

func (h *Handler) UpdateEvent(
	w http.ResponseWriter,
	r *http.Request,
) {
	fmt.Println(w, r)
}

func (h *Handler) DeleteEvent(
	w http.ResponseWriter,
	r *http.Request,
) {
	fmt.Println(w, r)
}

func (h *Handler) ListEvent(
	w http.ResponseWriter,
	r *http.Request,
) {
	fmt.Println(w, r)
}
