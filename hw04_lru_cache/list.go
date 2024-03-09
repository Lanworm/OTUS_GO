package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length       int
	firstElement *ListItem
	lastElement  *ListItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.firstElement
}

func (l *list) Back() *ListItem {
	return l.lastElement
}

func (l *list) PushFront(v interface{}) *ListItem {
	tmp := &ListItem{Value: v, Next: l.firstElement, Prev: nil}

	if l.Len() == 0 {
		l.firstElement = tmp
		l.lastElement = tmp
		l.length++
		return l.firstElement
	}

	if l.Len() != 0 {
		l.firstElement.Prev = tmp
		l.firstElement = tmp
		l.length++
	}
	return l.firstElement
}

func (l *list) PushBack(v interface{}) *ListItem {
	tmp := &ListItem{Value: v, Next: nil, Prev: l.lastElement}

	if l.Len() == 0 {
		l.firstElement = tmp
		l.lastElement = tmp
		l.length++
		return l.lastElement
	}

	if l.Len() != 0 {
		l.lastElement.Next = tmp
		l.lastElement = tmp
		l.length++
	}
	return l.lastElement
}

func (l *list) Remove(i *ListItem) {
	switch {
	case l.length == 1:
		zero := &ListItem{Value: nil, Next: nil, Prev: nil}
		l.firstElement = zero
		l.lastElement = zero
	case i == l.Back():
		l.lastElement = i.Prev
		i.Prev.Next = nil
	case i == l.Front():
		l.firstElement = i.Next
		i.Next.Prev = nil
	default:
		tmp := i.Prev.Next
		i.Prev.Next = i.Next
		i.Next.Prev = tmp
	}
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	switch {
	case i.Prev == nil:
		return
	case i.Next == nil:
		l.lastElement = l.lastElement.Prev
		l.lastElement.Next = nil
	default:
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	i.Next = l.firstElement
	i.Prev = nil
	l.firstElement = i
}

func NewList() List {
	return new(list)
}
