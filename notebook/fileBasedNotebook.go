package notebook

import (
	"todo-cli/notebook/internal/persistence"
)

const (
	defaultFilename = ".notebook"
)

// Note: concurrent update is not supported yet.
type fileBasedNotebook struct {
	protectedNotebook
	notebookFile persistence.JsonFile
}

func (s *fileBasedNotebook) Add(summary string) (Item, error) {
	var items []Item
	if err := s.notebookFile.Read(&items); err != nil {
		return Item{}, err
	}
	newItem := Item{
		// Assuming no items can be deleted, then `len(items)` is the max Id.
		Id:      len(items) + 1,
		Summary: summary,
		Done:    false,
	}
	if err := s.notebookFile.Write(append(items, newItem)); err != nil {
		return Item{}, err
	}
	return newItem, nil
}

func (s *fileBasedNotebook) Done(id ItemId) error {
	var items []Item
	if err := s.notebookFile.Read(&items); err != nil {
		return err
	}
	for i := range items {
		if items[i].Id == id {
			items[i].Done = true
		}
	}
	return s.notebookFile.Write(items)
}

func (s *fileBasedNotebook) ListUndoneItems() ([]Item, error) {
	var items []Item
	if err := s.notebookFile.Read(&items); err != nil {
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

func (s *fileBasedNotebook) ListAllItems() (items []Item, doneItems int, err error) {
	if err := s.notebookFile.Read(&items); err != nil {
		return nil, 0, err
	}
	for _, item := range items {
		if item.Done {
			doneItems++
		}
	}
	return items, doneItems, nil
}
