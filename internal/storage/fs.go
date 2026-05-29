package storage

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"

	"github.com/nskforward/gate4/pkg/tools"
)

type FileObjectStorage struct {
	root *os.Root
}

func NewFileObjectStorage(rootDir string) *FileObjectStorage {
	normalized := tools.Path(rootDir)
	err := os.MkdirAll(normalized, os.ModePerm)
	if err != nil {
		slog.Error("storage.fs: cannot create the root dir", "path", normalized, "reason", err.Error())
		os.Exit(1)
	}
	root, err := os.OpenRoot(normalized)
	if err != nil {
		slog.Error("storage.fs: cannot open the root dir", "path", normalized, "reason", err.Error())
		os.Exit(1)
	}
	return &FileObjectStorage{
		root: root,
	}
}

func (fs *FileObjectStorage) Keys(ctx context.Context) []string {
	files, err := os.ReadDir(fs.root.Name())
	if err != nil {
		log.Fatal(err)
	}
	keys := make([]string, 0, 32)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		keys = append(keys, file.Name())
	}
	return keys
}

func (fs *FileObjectStorage) Save(ctx context.Context, key string, src any) error {
	if key == "" {
		return ErrKeyEmpty
	}
	f, err := fs.root.Create(key)
	if err != nil {
		return wrapFileErr(err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(src)
}

func (fs *FileObjectStorage) Read(ctx context.Context, key string, dst any) error {
	if key == "" {
		return ErrKeyEmpty
	}
	f, err := fs.root.Open(key)
	if err != nil {
		return wrapFileErr(err)
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(dst)
}

func (fs *FileObjectStorage) Delete(ctx context.Context, key string) error {
	if key == "" {
		return ErrKeyEmpty
	}
	err := fs.root.Remove(key)
	return wrapFileErr(err)
}

func wrapFileErr(err error) error {
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return ErrNotFound
	}
	return err
}
