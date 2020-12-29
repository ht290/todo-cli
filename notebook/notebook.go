package notebook

// Keep all public interfaces in this file.

type ItemId = int

type Item struct {
	Id      ItemId
	Summary string
	Done    bool
}

type Items interface {
	Add(summary string) (Item, error)
	Done(itemId ItemId) error
	ListAllItems() (items []Item, doneItems int, err error)
	ListUndoneItems() ([]Item, error)
}

func NewForSingleUser() *singleUserNotebook {
	return &singleUserNotebook{filename: defaultFilename}
}
