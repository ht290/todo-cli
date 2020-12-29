package notebook

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	defaultFilename = ".notebook"
)

// Note: concurrent update is not supported yet.
type singleUserNotebook struct {
	filename string
}

func (s *singleUserNotebook) Add(summary string) (Item, error) {
	items, err := s.loadItemsFromFile()
	if err != nil {
		return Item{}, err
	}
	newItem := Item{
		// Assuming no items can be deleted, then `len(items)` is the max Id.
		Id:      len(items) + 1,
		Summary: summary,
		Done:    false,
	}
	if err := s.writeItemsToFile(append(items, newItem)); err != nil {
		return Item{}, err
	}
	return newItem, nil
}

func (s *singleUserNotebook) Done(id ItemId) error {
	items, err := s.loadItemsFromFile()
	if err != nil {
		return err
	}
	for i := range items {
		if items[i].Id == id {
			items[i].Done = true
		}
	}
	return s.writeItemsToFile(items)
}

func (s *singleUserNotebook) ListUndoneItems() ([]Item, error) {
	items, err := s.loadItemsFromFile()
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

func (s *singleUserNotebook) ListAllItems() (items []Item, doneItems int, err error) {
	items, err = s.loadItemsFromFile()
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

func (s *singleUserNotebook) writeItemsToFile(items []Item) error {
	serializedItems, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(s.filename, serializedItems, 0644); err != nil {
		return err
	}
	return nil
}

func (s *singleUserNotebook) loadItemsFromFile() ([]Item, error) {
	if !s.fileExists() {
		return nil, nil
	}
	serializedItems, err := ioutil.ReadFile(s.filename)
	if err != nil {
		return nil, err
	}
	if len(serializedItems) == 0 {
		return nil, nil
	}
	var items []Item
	if err := json.Unmarshal(serializedItems, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *singleUserNotebook) fileExists() bool {
	_, err := os.Stat(s.filename)
	return err == nil
}
