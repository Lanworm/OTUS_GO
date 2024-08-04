package httphandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/enum"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/server/http/dto"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/service"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage"
)

type Handler struct {
	logger  *logger.Logger
	service *service.Event
}

func NewHandler(
	logger *logger.Logger,
	service *service.Event,
) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) CreateEvent(
	w http.ResponseWriter,
	r *http.Request,
) {
	body, err := readBody(r)
	if err != nil {
		writeError(http.StatusInternalServerError, w, err.Error())
		return
	}

	var createReq dto.CreateEvent

	err = json.Unmarshal(body, &createReq)
	if err != nil {
		writeError(http.StatusBadRequest, w, err.Error())
		return
	}

	id, err := h.service.CreateEvent(&storage.Event{
		Title:         createReq.Title,
		StartDatetime: createReq.StartTime,
		EndDatetime:   createReq.EndTime,
		Description:   createReq.Description,
		UserID:        createReq.UserID,
		RemindBefore:  createReq.RemindBefore,
	})
	if err != nil {
		writeError(http.StatusInternalServerError, w, err.Error())
		return
	}

	writeResult(http.StatusOK, w, dto.CreateReply{
		ID: id,
	})
}

func (h *Handler) GetEvent(
	w http.ResponseWriter,
	r *http.Request,
) {
	eventID := strings.TrimSpace(r.URL.Query().Get("event_id"))
	if eventID == "" {
		writeError(http.StatusBadRequest, w, "parameter event_id is invalid")
		return
	}

	evt, err := h.service.GetEvent(eventID)
	if err != nil {
		writeError(http.StatusInternalServerError, w, err.Error())
		return
	}

	writeResult(http.StatusOK, w, evt)
}

func (h *Handler) UpdateEvent(
	w http.ResponseWriter,
	r *http.Request,
) {
	body, err := readBody(r)
	if err != nil {
		writeError(http.StatusInternalServerError, w, err.Error())
		return
	}

	var updateReq dto.UpdateEvent

	err = json.Unmarshal(body, &updateReq)
	if err != nil {
		writeError(http.StatusBadRequest, w, err.Error())
		return
	}

	if strings.TrimSpace(updateReq.ID) == "" {
		writeError(http.StatusBadRequest, w, "parameter event.id is invalid")
		return
	}

	err = h.service.UpdateEvent(&storage.Event{
		ID:            updateReq.ID,
		Title:         updateReq.Title,
		StartDatetime: updateReq.StartTime,
		EndDatetime:   updateReq.EndTime,
		Description:   updateReq.Description,
		RemindBefore:  updateReq.RemindBefore,
	})
	if err != nil {
		writeError(http.StatusInternalServerError, w, err.Error())
		return
	}
}

func (h *Handler) DeleteEvent(
	w http.ResponseWriter,
	r *http.Request,
) {
	eventID := strings.TrimSpace(r.URL.Query().Get("event_id"))
	if eventID == "" {
		writeError(http.StatusBadRequest, w, "parameter event_id is invalid")
		return
	}

	err := h.service.DeleteEvent(eventID)
	if err != nil {
		writeError(http.StatusInternalServerError, w, err.Error())
		return
	}
}

func (h *Handler) ListEvent(
	w http.ResponseWriter,
	r *http.Request,
) {
	dateRange, err := enum.NewRangeDurationByString(r.URL.Query().Get("range"))
	if err != nil {
		writeError(http.StatusBadRequest, w, err.Error())
	}

	list, err := h.service.ListEvent(dateRange)
	if err != nil {
		writeError(http.StatusInternalServerError, w, err.Error())
		return
	}

	js, err := json.Marshal(dto.ListReply{Result: list})
	if err != nil {
		writeError(http.StatusInternalServerError, w, err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}
}

func readBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("read input body: %w", err)
	}
	defer r.Body.Close()

	return body, nil
}

func writeResult(statusCode int, w http.ResponseWriter, reply interface{}) {
	js, err := json.Marshal(reply)
	if err != nil {
		writeError(http.StatusInternalServerError, w, err.Error())
	} else {
		w.WriteHeader(statusCode)
		w.Write(js)
	}
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

func (h *Handler) Health(
	w http.ResponseWriter,
	_ *http.Request,
) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Ready(
	w http.ResponseWriter,
	_ *http.Request,
) {
	w.WriteHeader(http.StatusOK)
}
