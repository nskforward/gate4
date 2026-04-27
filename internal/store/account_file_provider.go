package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nskforward/gate4/pkg/pb"
)

type AccountFileProvider struct {
	file string
}

func NewAccountFileProvider(rootDir string) *AccountFileProvider {
	return &AccountFileProvider{
		file: filepath.Join(rootDir, "accounts.json"),
	}
}

func (s *AccountFileProvider) Load() (map[string]*pb.Account, error) {
	f, err := os.Open(s.file)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]*pb.Account), nil
		}
		return nil, fmt.Errorf("cannot load accounts from file: %w", err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	items := make(map[string]*pb.Account)
	err = decoder.Decode(&items)
	if err != nil {
		return nil, fmt.Errorf("cannot decode accounts from file: %w", err)
	}
	return items, nil
}

func (s *AccountFileProvider) Save(items map[string]*pb.Account) error {
	f, err := os.Create(s.file)
	if err != nil {
		return fmt.Errorf("cannot open accounts file for writing: %w", err)
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	err = encoder.Encode(items)
	if err != nil {
		return fmt.Errorf("cannot encode accounts for writing: %w", err)
	}
	return nil
}
