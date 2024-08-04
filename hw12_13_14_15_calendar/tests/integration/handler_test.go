package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/server/http/dto"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/server/http/httphandler"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/service"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateEvent(t *testing.T) {
	logBytes := make([]byte, 0, 1000)
	bLog := bytes.NewBuffer(logBytes)

	logg, err := logger.New("DEBUG", bLog)
	assert.NoErrorf(t, err, "fail initialize logger")

	stor := memorystorage.New()

	s := service.NewEventService(logg, stor)
	handl := httphandler.NewHandler(logg, s)

	srv := http.HandlerFunc(handl.CreateEvent)

	bodyStr := bytes.NewBufferString(`{
		"title": "1122",
		"description": "description",
		"user_id": "f8e51dcd-f8fd-459c-81c3-2c873e40d747",
		"remind_before": 1023,
		"start_time": "2024-06-04T18:42:33.177Z",
		"end_time": "2024-06-04T20:42:33.177Z"
	}`)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/event/create",
		bodyStr,
	)

	srv.ServeHTTP(rec, req)

	result := rec.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
	defer result.Body.Close()

	var createReply dto.CreateReply
	err = json.Unmarshal(rec.Body.Bytes(), &createReply)

	assert.NoError(t, err, "unmarshal response")
	assert.True(t, createReply.ID != "")
}

func TestHandler_DeleteEvent(t *testing.T) {
	logBytes := make([]byte, 0, 1000)
	bLog := bytes.NewBuffer(logBytes)

	logg, err := logger.New("DEBUG", bLog)
	assert.NoErrorf(t, err, "fail initialize logger")

	stor := memorystorage.New()
	delEvent, err := stor.Add(&storage.Event{
		Title:         "111",
		StartDatetime: time.Time{},
		EndDatetime:   time.Time{},
		Description:   "",
		UserID:        "",
		RemindBefore:  0,
	})
	assert.NoErrorf(t, err, "create event storage")

	s := service.NewEventService(logg, stor)
	handl := httphandler.NewHandler(logg, s)

	srv := http.HandlerFunc(handl.DeleteEvent)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("/event/delete?event_id=%s", delEvent),
		nil,
	)
	assert.NoErrorf(t, err, "create delete request")

	srv.ServeHTTP(rec, req)

	result := rec.Result()
	defer result.Body.Close()
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestHandler_GetEvent(t *testing.T) {
	logBytes := make([]byte, 0, 1000)
	bLog := bytes.NewBuffer(logBytes)

	logg, err := logger.New("DEBUG", bLog)
	assert.NoErrorf(t, err, "fail initialize logger")

	stor := memorystorage.New()
	getEvent, err := stor.Add(&storage.Event{
		Title:         "111",
		StartDatetime: time.Time{},
		EndDatetime:   time.Time{},
		Description:   "",
		UserID:        "",
		RemindBefore:  0,
	})
	assert.NoErrorf(t, err, "create event storage")

	s := service.NewEventService(logg, stor)
	handl := httphandler.NewHandler(logg, s)

	srv := http.HandlerFunc(handl.GetEvent)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/event?event_id=%s", getEvent),
		nil,
	)
	assert.NoErrorf(t, err, "create get event request")

	srv.ServeHTTP(rec, req)

	result := rec.Result()
	defer result.Body.Close()
	assert.Equal(t, http.StatusOK, result.StatusCode)

	var res storage.Event
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoErrorf(t, err, "unmarshal body")

	assert.Equal(t, getEvent, res.ID)
	assert.Equal(t, "111", res.Title)
}

func TestHandler_ListEvent(t *testing.T) {
	logBytes := make([]byte, 0, 1000)
	bLog := bytes.NewBuffer(logBytes)

	logg, err := logger.New("DEBUG", bLog)
	assert.NoErrorf(t, err, "fail initialize logger")

	stor := memorystorage.New()
	evt1, err := stor.Add(&storage.Event{
		Title:         "111",
		StartDatetime: time.Now(),
		EndDatetime:   time.Time{},
		Description:   "",
		UserID:        "",
		RemindBefore:  0,
	})
	assert.NoErrorf(t, err, "create event storage 1")

	evt2, err := stor.Add(&storage.Event{
		Title:         "222",
		StartDatetime: time.Now(),
		EndDatetime:   time.Time{},
		Description:   "",
		UserID:        "",
		RemindBefore:  0,
	})
	assert.NoErrorf(t, err, "create event storage 2")

	s := service.NewEventService(logg, stor)
	handl := httphandler.NewHandler(logg, s)

	srv := http.HandlerFunc(handl.ListEvent)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodGet,
		"/event/list?range=day",
		nil,
	)
	assert.NoErrorf(t, err, "create list event request")

	srv.ServeHTTP(rec, req)

	result := rec.Result()
	defer result.Body.Close()
	assert.Equal(t, http.StatusOK, result.StatusCode)

	var res dto.ListReply
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoErrorf(t, err, "unmarshal body")

	assert.Equal(t, 2, len(res.Result))
	mapping := map[string]string{
		evt1: "111",
		evt2: "222",
	}
	for _, event := range res.Result {
		if evt, ok := mapping[event.ID]; !ok {
			t.Errorf("unexpected id=%s", event.ID)
		} else {
			assert.Equal(t, evt, event.Title)
		}
	}
}

func TestHandler_UpdateEvent(t *testing.T) {
	logBytes := make([]byte, 0, 1000)
	bLog := bytes.NewBuffer(logBytes)

	logg, err := logger.New("DEBUG", bLog)
	assert.NoErrorf(t, err, "fail initialize logger")

	stor := memorystorage.New()
	updateEvt, err := stor.Add(&storage.Event{
		Title:         "111",
		StartDatetime: time.Now(),
		EndDatetime:   time.Now(),
		Description:   "",
		UserID:        "",
		RemindBefore:  0,
	})
	assert.NoErrorf(t, err, "create event storage")

	s := service.NewEventService(logg, stor)
	handl := httphandler.NewHandler(logg, s)

	srv := http.HandlerFunc(handl.UpdateEvent)

	reqEvent, err := json.Marshal(&storage.Event{
		ID:            updateEvt,
		Title:         "111",
		StartDatetime: time.Now(),
		EndDatetime:   time.Now(),
		Description:   "111_description",
		UserID:        "",
		RemindBefore:  0,
	})
	assert.NoErrorf(t, err, "find item in storage")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/event?event_id=%s", updateEvt),
		bytes.NewBuffer(reqEvent),
	)
	assert.NoErrorf(t, err, "update event request")

	resEvt, err := stor.FindItem(updateEvt)
	assert.NoErrorf(t, err, "find event before update")
	assert.Equal(t, resEvt.Description, "")

	srv.ServeHTTP(rec, req)

	result := rec.Result()
	defer result.Body.Close()
	assert.Equal(t, http.StatusOK, result.StatusCode)

	resEvt, err = stor.FindItem(updateEvt)
	assert.NoErrorf(t, err, "find updated event")
	assert.Equal(t, resEvt.Description, "111_description")
}
