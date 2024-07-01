package httphandler

import (
	"bytes"
	"encoding/json"
	"image/jpeg"
	"net/http"

	"github.com/Lanworm/OTUS_GO/final_project/internal/http/server/dto"
	"github.com/Lanworm/OTUS_GO/final_project/internal/logger"
	"github.com/Lanworm/OTUS_GO/final_project/internal/service"
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
	imgParams, err := service.PrepareImgParams(r.URL)
	if err != nil {
		writeError(http.StatusInternalServerError, w, err.Error())
		h.logger.Error(err.Error())
		return
	}
	img, err := service.ResizeImg(imgParams)
	if err != nil {
		writeError(http.StatusInternalServerError, w, err.Error())
		h.logger.Error(err.Error())
		return
	}
	w.Header().Set("Content-Type", "image/png")
	buf := new(bytes.Buffer)
	encodeErr := jpeg.Encode(buf, img, nil)
	if encodeErr != nil {
		writeError(http.StatusInternalServerError, w, encodeErr.Error())
		h.logger.Error(encodeErr.Error())
		return
	}
	w.Write(buf.Bytes())
}

func writeError(
	statusCode int,
	w http.ResponseWriter,
	msg string,
) {
	js, err := json.Marshal(dto.Result{Message: msg})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(statusCode)
		w.Write(js)
	}
}
