package notebook

import (
	"todo-cli/notebook/internal/persistence"
)

type ItemStorage interface {
	list(user UserId) ([]Item, error)
	create(draft Item) (result Item, err error)
	update(id ItemId, isDone bool) error
}

type itemsJsonFile struct {
	file persistence.JsonFile
}

func InitItemsStorage(filename string) *itemsJsonFile {
	return &itemsJsonFile{
		file: persistence.JsonFile{
			Filename: filename,
		},
	}
}

func (i *itemsJsonFile) list(user UserId) ([]Item, error) {
	var allUsersItems []Item
	if err := i.file.Read(&allUsersItems); err != nil {
		return nil, err
	}
	var myItems []Item
	for _, item := range allUsersItems {
		if item.Author == user {
			myItems = append(myItems, item)
		}
	}
	return myItems, nil
}

func (i *itemsJsonFile) create(draft Item) (result Item, err error) {
	var items []Item
	if err := i.file.Read(&items); err != nil {
		return Item{}, err
	}
	maxId, err := i.findMaxId(draft.Author, items)
	if err != nil {
		return Item{}, err
	}
	draft.Id = maxId + 1
	return draft, i.file.Write(append(items, draft))
}

func (i *itemsJsonFile) update(id ItemId, isDone bool) error {
	var items []Item
	if err := i.file.Read(&items); err != nil {
		return err
	}
	for i := range items {
		if items[i].Id == id {
			items[i].Done = isDone
		}
	}
	return i.file.Write(items)
}

func (i *itemsJsonFile) findMaxId(user UserId, items []Item) (ItemId, error) {
	maxId := 0
	for _, item := range items {
		if item.Author == user && item.Id > maxId {
			maxId = item.Id
		}
	}
	return maxId, nil
}
