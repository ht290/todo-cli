package notebook

import (
	"fmt"
	"os"
	"path"
	"todo-cli/notebook/internal/persistence"
)

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

type Protected interface {
	Lock() error
	Unlock(username, password string) error
	Create(username, password string) error
}

func New() *fileBasedNotebook {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to get home dir. Use working dir instead.")
		homeDir = "."
	}
	return &fileBasedNotebook{
		protectedNotebook: newProtectedNotebook(path.Join(homeDir, defaultAccountFilename)),
		notebookFile:      persistence.JsonFile{Filename: defaultFilename},
	}
}
