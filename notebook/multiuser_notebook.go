package notebook

const (
	DefaultNotebookFilename = ".notebook"
)

// Note: concurrent update is not supported yet.
type multiuserNotebook struct {
	storage   ItemStorage
	encryptor ItemEncryptor
}

func (s *multiuserNotebook) newDraftItem(summary string) (Item, error) {
	return Item{
		Author:  s.encryptor.getActiveUser(),
		Summary: summary,
	}, nil
}

func (s *multiuserNotebook) Add(summary string) (Item, error) {
	draft, err := s.newDraftItem(summary)
	if err != nil {
		return Item{}, err
	}
	encrypted, err := s.encryptor.encrypt(draft)
	if err != nil {
		return Item{}, err
	}
	writtenItem, err := s.storage.create(encrypted)
	if err != nil {
		return Item{}, err
	}
	result, err := s.encryptor.decrypt(writtenItem)
	return result, err
}

func (s *multiuserNotebook) Done(id ItemId) error {
	return s.storage.update(id, true)
}

func (s *multiuserNotebook) ListUndoneItems() ([]Item, error) {
	items, _, err := s.ListAllItems()
	if err != nil {
		return nil, err
	}
	var undoneItems []Item
	for _, item := range items {
		if !item.Done {
			undoneItems = append(undoneItems, item)
		}
	}
	return undoneItems, nil
}

func (s *multiuserNotebook) ListAllItems() (items []Item, doneItems int, err error) {
	encryptedItems, err := s.storage.list(s.encryptor.getActiveUser())
	if err != nil {
		return nil, 0, err
	}
	items, err = s.encryptor.decryptAll(encryptedItems)
	if err != nil {
		return nil, 0, err
	}
	for _, item := range items {
		if item.Done {
			doneItems++
		}
	}
	return items, doneItems, nil
}

func (s *multiuserNotebook) Create(username, password string) error {
	return s.encryptor.createKey(username, password)
}

func (s *multiuserNotebook) Unlock(username, password string) error {
	return s.encryptor.retrieveKey(username, password)
}

func (s *multiuserNotebook) Lock() error {
	return s.encryptor.returnKey()
}
