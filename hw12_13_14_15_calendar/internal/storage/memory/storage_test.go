package memorystorage

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestConcurrencyStorage(t *testing.T) {
	testFunc := func() {
		strg := New()

		pipe := make(chan int, 3)

		wg := sync.WaitGroup{}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for item := range pipe {
					_, err := strg.Add(&storage.Event{Title: strconv.Itoa(item)})
					assert.NoError(t, err)
				}
			}()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				pipe <- i
			}

			close(pipe)
		}()

		wg.Wait()
	}

	for i := 0; i < 100; i++ {
		assert.NotPanics(t, testFunc)
	}
}

func TestStorage_AddAndFind(t *testing.T) {
	strg := New()
	id, err := strg.Add(&storage.Event{
		ID:            "123",
		Title:         "123",
		StartDatetime: time.Time{},
		EndDatetime:   time.Time{},
		Description:   "123",
		UserID:        "123",
		RemindBefore:  123,
	})
	assert.NoError(t, err)

	item, err := strg.FindItem(id)
	assert.NoError(t, err)

	_, err = strg.FindItem("123")
	assert.ErrorIs(t, storage.ErrEventNotFound, err)

	assert.Equalf(
		t,
		id,
		item.ID,
		"item id %v, is not equal id %v from return Add method",
		item.ID,
		id,
	)
}

func TestStorage_List(t *testing.T) {
	strg := New()
	expected := getEvents()

	m := make(map[string]time.Time, len(expected))
	for i, ev := range expected {
		stDatetime := time.Now().Add((time.Hour * 24) * time.Duration(i+1))
		ev.StartDatetime = stDatetime
		id, err := strg.Add(&ev)

		m[id] = stDatetime

		assert.NoError(t, err)
	}

	actual := make([]storage.Event, 0, len(expected))
	for id, startTime := range m {
		endTime := startTime.Add(time.Hour + 1)

		items, err := strg.ListRange(&startTime, &endTime)
		assert.NoError(t, err)
		assert.Len(t, items, 1)

		for _, item := range items {
			item := item
			if item.ID == id {
				actual = append(actual, item)
			} else {
				t.Errorf("ListRange not returned element by id %s", id)
			}
		}
	}

	assert.True(t, len(actual) == len(expected))
}

func TestStorage_Delete(t *testing.T) {
	strg := New()
	item1, err := strg.Add(&storage.Event{
		ID:            "1",
		Title:         "1",
		StartDatetime: time.Time{},
		EndDatetime:   time.Time{},
		Description:   "1",
		UserID:        "1",
		RemindBefore:  int64(1),
	})
	assert.NoError(t, err)

	item2, err := strg.Add(&storage.Event{
		ID:            "2",
		Title:         "2",
		StartDatetime: time.Time{},
		EndDatetime:   time.Time{},
		Description:   "2",
		UserID:        "2",
		RemindBefore:  int64(2),
	})
	assert.NoError(t, err)

	err = strg.Delete(item2)
	assert.NoError(t, err)

	_, err = strg.FindItem(item1)
	assert.NoErrorf(t, err, "find item1 return error")

	_, err = strg.FindItem(item2)
	assert.ErrorIs(t, storage.ErrEventNotFound, err, "find item2  return error")
}

func getEvents() []storage.Event {
	return []storage.Event{
		{
			ID:            "1",
			Title:         "1",
			StartDatetime: time.Time{},
			EndDatetime:   time.Time{},
			Description:   "1",
			UserID:        "1",
			RemindBefore:  int64(1),
		},
		{
			ID:            "2",
			Title:         "2",
			StartDatetime: time.Time{},
			EndDatetime:   time.Time{},
			Description:   "2",
			UserID:        "2",
			RemindBefore:  int64(2),
		},
		{
			ID:            "3",
			Title:         "3",
			StartDatetime: time.Time{},
			EndDatetime:   time.Time{},
			Description:   "3",
			UserID:        "3",
			RemindBefore:  int64(3),
		},
		{
			ID:            "4",
			Title:         "4",
			StartDatetime: time.Time{},
			EndDatetime:   time.Time{},
			Description:   "4",
			UserID:        "4",
			RemindBefore:  int64(4),
		},
		{
			ID:            "5",
			Title:         "5",
			StartDatetime: time.Time{},
			EndDatetime:   time.Time{},
			Description:   "5",
			UserID:        "5",
			RemindBefore:  int64(5),
		},
	}
}
