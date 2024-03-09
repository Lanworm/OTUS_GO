package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type keyVal struct {
	k Key
	v interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	if l.queue.Len() == l.capacity {
		delete(l.items, l.queue.Back().Value.(*keyVal).k)
		l.queue.Remove(l.queue.Back())
	}
	v, ok := l.items[key]
	if ok {
		l.items[key] = &ListItem{Value: value}
		l.queue.MoveToFront(v)
	} else {
		l.items[key] = &ListItem{Value: value}
		l.queue.PushFront(&keyVal{key, value})
	}
	return ok
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	v, ok := l.items[key]
	if !ok {
		return nil, ok
	}
	l.queue.MoveToFront(v)
	return l.items[key].Value, ok
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem)
	l.capacity = 0
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
