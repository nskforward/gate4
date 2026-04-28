package object

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileStorage[T any] struct {
	filename string
	items    map[string]T
	mx       sync.RWMutex
}

func NewFileStorage[T any](filename string) (*FileStorage[T], error) {
	normalized, err := normalizeFilename(filename)
	if err != nil {
		return nil, fmt.Errorf("bad filename: %w", err)
	}
	s := &FileStorage[T]{
		filename: normalized,
		items:    make(map[string]T),
	}
	err = s.readFromFile()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *FileStorage[T]) Set(key string, obj T) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.items[key] = obj
	return s.saveToFile()
}

func (s *FileStorage[T]) Get(key string) (T, bool) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	obj, ok := s.items[key]
	return obj, ok
}

func (s *FileStorage[T]) Del(key string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.items, key)
	return s.saveToFile()
}

func (s *FileStorage[T]) Keys() []string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	result := make([]string, 0, len(s.items))
	for key := range s.items {
		result = append(result, key)
	}
	return result
}

func (s *FileStorage[T]) Objects() []T {
	s.mx.RLock()
	defer s.mx.RUnlock()
	result := make([]T, 0, len(s.items))
	for _, obj := range s.items {
		result = append(result, obj)
	}
	return result
}

func (s *FileStorage[T]) readFromFile() error {
	f, err := os.Open(s.filename)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&s.items)
}

func (s *FileStorage[T]) saveToFile() error {
	f, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(s.items)
}

func normalizeFilename(in string) (string, error) {
	if in == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}
	if in[0] == '/' {
		return in, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, in), nil
}
