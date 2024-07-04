package lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type CacheListItem struct {
	value interface{}
	key   Key
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	cacheListItem, ok := c.items[key]

	if ok {
		cacheItem := cacheListItem.Value.(CacheListItem)
		cacheItem.value = value
		cacheListItem.Value = cacheItem

		c.queue.MoveToFront(cacheListItem)

		return true
	}

	newCacheItem := CacheListItem{
		value: value,
		key:   key,
	}

	if c.queue.Len() >= c.capacity {
		lastListItem := c.queue.Back()
		cad := lastListItem.Value.(CacheListItem)
		delete(c.items, cad.key)
	}

	c.items[key] = c.queue.PushFront(newCacheItem)

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	cacheItem, ok := c.items[key]

	if ok {
		c.queue.MoveToFront(cacheItem)
		cad := cacheItem.Value.(CacheListItem)

		return cad.value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
