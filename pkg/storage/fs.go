package storage

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"

	"github.com/nskforward/gate4/pkg/common"
)

type FileStorage struct {
	root *os.Root
}

func NewFileStorage(rootDir string) *FileStorage {
	normalized := common.Path(rootDir)
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
	return &FileStorage{
		root: root,
	}
}

func (fs *FileStorage) Keys(ctx context.Context) []string {
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

func (fs *FileStorage) Save(ctx context.Context, key string, src any) error {
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

func (fs *FileStorage) Find(ctx context.Context, key string, dst any) error {
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

func (fs *FileStorage) Delete(ctx context.Context, key string) error {
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
