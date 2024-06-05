package storage

type IStorage interface {
	Add(item Event) (string, error)
	Update(id string, item Event) error
	Delete(id string) error
	FindItem(id string) (Event, error)
	List() ([]Event, error)
}
