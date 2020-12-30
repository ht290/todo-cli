package persistence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ReaderWriter interface {
	Read(out interface{}) error
	Write(value interface{}) error
}

type JsonFile struct {
	Filename string
}

func (j *JsonFile) Read(out interface{}) error {
	if !j.fileExists() {
		return nil
	}
	serializedItems, err := ioutil.ReadFile(j.Filename)
	if err != nil {
		return err
	}
	if len(serializedItems) == 0 {
		return nil
	}
	return json.Unmarshal(serializedItems, out)
}

func (j *JsonFile) fileExists() bool {
	_, err := os.Stat(j.Filename)
	return err == nil
}

func (j *JsonFile) Write(value interface{}) error {
	serializedItems, err := json.MarshalIndent(value, "", " ")
	if err != nil {
		return err
	}
	fmt.Printf("s %+v", serializedItems)
	return ioutil.WriteFile(j.Filename, serializedItems, 0644)
}
