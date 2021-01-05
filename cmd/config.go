package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"todo-cli/notebook"
)

const (
	credentialFileKey = "CredentialFile"
	notebookFileKey   = "NotebookFile"
)

func initNotebook() notebook.Notebook {
	storage := notebook.InitItemsStorage(viper.GetString(notebookFileKey))
	newNotebook, err := notebook.New(initItemEncryptor(), storage)
	if err != nil {
		log.Fatal(err)
	}
	return newNotebook
}

func initItemEncryptor() notebook.ItemEncryptor {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to get home dir. Use working dir instead.")
		homeDir = "."
	}
	vault, err := notebook.UseLocalItemEncryptor(path.Join(homeDir, viper.GetString(credentialFileKey)))
	if err != nil {
		log.Fatal(err)
	}
	return vault
}

func init() {
	viper.SetDefault(credentialFileKey, notebook.DefaultCredentialFilename)
	viper.SetDefault(notebookFileKey, notebook.DefaultNotebookFilename)
}
