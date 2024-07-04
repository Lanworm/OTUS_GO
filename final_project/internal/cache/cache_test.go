package lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("simple ejection", func(t *testing.T) {
		c := NewCache(3)
		c.Set("a1", 1)
		c.Set("a2", 2)
		c.Set("a3", 3)
		c.Set("a4", 4)

		val, ok := c.Get("a1")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("ejection old element", func(t *testing.T) {
		c := NewCache(3)

		c.Set("a1", 1)
		c.Set("a2", 2)
		c.Set("a3", 3)

		val, ok := c.Get("a3")
		require.True(t, ok)
		require.NotNil(t, val)

		c.Set("a2", 40)

		ok = c.Set("a4", 104)
		require.False(t, ok)

		val, ok = c.Get("a1")
		require.False(t, ok)
		require.Nil(t, val)
	})
}
