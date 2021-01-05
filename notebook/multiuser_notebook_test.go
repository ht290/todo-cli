package notebook

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

// Due to time constraint, only adding an integration test here.
func TestCRUD(t *testing.T) {
	notebookFile, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	defer os.Remove(notebookFile.Name())
	credentialFile, err := ioutil.TempFile("", "cred")
	require.NoError(t, err)
	defer os.Remove(credentialFile.Name())

	vault, err := UseLocalItemEncryptor(credentialFile.Name())
	require.NoError(t, err)
	notebook, err := New(vault, InitItemsStorage(notebookFile.Name()))

	t.Run("Add a user and login", func(t *testing.T) {
		const (
			username = "jim"
			password = "pwd123"
		)
		require.NoError(t, notebook.Create(username, password))
		require.NoError(t, notebook.Unlock(username, password))
	})

	t.Run("List an empty notebook", func(t *testing.T) {
		items, doneItems, err := notebook.ListAllItems()
		require.NoError(t, err)
		require.Equal(t, 0, doneItems)
		require.Nil(t, items)
	})

	t.Run("Add some items", func(t *testing.T) {
		_, err := notebook.Add("call mum")
		require.NoError(t, err)
		_, err = notebook.Add("call dad")
		require.NoError(t, err)
		items, doneItems, err := notebook.ListAllItems()
		require.NoError(t, err)
		require.Equal(t, 0, doneItems)
		require.Equal(t, []Item{{
			Author:  "jim",
			Id:      1,
			Summary: "call mum",
			Done:    false,
		}, {
			Author:  "jim",
			Id:      2,
			Summary: "call dad",
			Done:    false,
		},
		}, items)
	})

	t.Run("Do some items", func(t *testing.T) {
		require.NoError(t, notebook.Done(2))
		items, doneItems, err := notebook.ListAllItems()
		require.NoError(t, err)
		require.Equal(t, 1, doneItems)
		require.Equal(t, []Item{{
			Author:  "jim",
			Id:      1,
			Summary: "call mum",
			Done:    false,
		}, {
			Author:  "jim",
			Id:      2,
			Summary: "call dad",
			Done:    true,
		},
		}, items)
	})
}
