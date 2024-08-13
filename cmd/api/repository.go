package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

type CounterRepository interface {
	Save(requests map[int64]int) error
	Get() (map[int64]int, error)
}

type FileStorage struct {
	filename string
}

func NewFileStorage(fn string) *FileStorage {
	return &FileStorage{filename: fn}
}

func (fs *FileStorage) Save(requests map[int64]int) error {
	file, err := os.OpenFile(fs.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error openning file: %v", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(requests)
	if err != nil {
		return fmt.Errorf("error encoding file: %v", err)
	}

	return nil
}

func (fs *FileStorage) Get() (map[int64]int, error) {
	file, err := os.Open(fs.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[int64]int), nil
		}
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var reqs map[int64]int
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&reqs)
	if err != nil {
		return nil, fmt.Errorf("error decoding data: %v", err)
	}

	return reqs, nil
}
