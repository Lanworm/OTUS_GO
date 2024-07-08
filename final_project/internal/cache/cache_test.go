package lrucache

import (
	"github.com/stretchr/testify/require"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"testing"
)

func TestLRUCache(t *testing.T) {
	cache := NewCache(2)

	// Проверка добавления и получения изображения из кеша
	img1 := image.NewRGBA(image.Rect(0, 0, 100, 100))
	cache.Set(Key("image1"), img1)

	retrievedImg1, found1 := cache.Get(Key("image1"))
	require.True(t, found1, "Изображение 'image1' не найдено в кеше")
	require.Equal(t, img1, retrievedImg1, "Изображение 'image1' не соответствует ожидаемому")

	// Проверка замещения изображения в кеше
	img2 := image.NewRGBA(image.Rect(0, 0, 200, 200))
	cache.Set(Key("image2"), img2)

	img3 := image.NewRGBA(image.Rect(0, 0, 300, 300))
	cache.Set(Key("image3"), img3)

	_, found2 := cache.Get(Key("image1"))
	require.False(t, found2, "Изображение 'image1' должно быть замещено")

	// Проверка очистки кеша
	cache.Clear()
	require.Empty(t, cache.(*lruCache).items, "Элементы кеша не были очищены")
	require.Zero(t, cache.(*lruCache).queue.Len(), "Длина очереди кеша должна быть нулевой")
}

func TestInitCache(t *testing.T) {
	capacity := 2
	testCache := NewCache(capacity)
	path := "../../test_images"

	err := InitCache(path, capacity, testCache)
	require.NoError(t, err, "Ошибка при инициализации кеша изображений")

	// Проверка добавления изображений в кеш
	retrievedImg1, found1 := testCache.Get(Key("gofer.jpg"))
	require.True(t, found1, "Изображение 'gofer.jpg' не найдено в кеше")
	require.NotNil(t, retrievedImg1, "Изображение 'gofer.jpg' не было добавлено в кеш")
}
