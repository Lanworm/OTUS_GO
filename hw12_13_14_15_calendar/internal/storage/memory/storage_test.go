package memorystorage

import (
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestConcurrencyStorage(t *testing.T) {
	storageExpected := New()

	for i := 1; i < 10; i++ {
		e := storage.Event{ID: strconv.Itoa(i)}
		_, err := storageExpected.Add(e)
		assert.NoError(t, err)
	}

	testStorage := New()
	wg := sync.WaitGroup{}

	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			e := storage.Event{ID: strconv.Itoa(i)}
			_, err := testStorage.Add(e)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	expected, err := storageExpected.List()
	assert.NoError(t, err)

	actual, err := testStorage.List()
	assert.NoError(t, err)

	sortFunc := func(a, b storage.Event) int {
		if a.ID > b.ID {
			return 1
		}

		if a.ID < b.ID {
			return -1
		}

		return 0
	}

	slices.SortFunc(expected, sortFunc)
	slices.SortFunc(actual, sortFunc)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("wanted %+v\n\rgot %+v", expected, actual)
	}
}
