package notebook

type ItemId = int

type UserId = string

type Item struct {
	Author UserId
	Id     ItemId
	// Plaintext summary must not be serialized and stored.
	Summary          string `json:"-"`
	EncryptedSummary []byte
	Done             bool
}

type Items interface {
	Add(summary string) (Item, error)
	Done(itemId ItemId) error
	ListAllItems() (items []Item, doneItems int, err error)
	ListUndoneItems() ([]Item, error)
}
