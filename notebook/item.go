package notebook

type ItemId = int

type UserId = string

type Item struct {
	Author           UserId
	Id               ItemId
	Summary          string
	EncryptedSummary []byte
	Done             bool
}

type Items interface {
	Add(summary string) (Item, error)
	Done(itemId ItemId) error
	ListAllItems() (items []Item, doneItems int, err error)
	ListUndoneItems() ([]Item, error)
}
