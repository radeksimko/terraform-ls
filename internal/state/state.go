package state

import (
	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/terraform-ls/internal/terraform/module"
)

var dbSchema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		providerSchemaTableName: providerSchemaTable,
	},
}

type StateStore struct {
	db        *memdb.MemDB
	modFinder module.ModuleFinder
}

func NewStateStore() (*StateStore, error) {
	db, err := memdb.NewMemDB(dbSchema)
	if err != nil {
		return nil, err
	}

	return &StateStore{
		db: db,
	}, nil
}

func (s *StateStore) SetModuleFinder(mf module.ModuleFinder) {
	s.modFinder = mf
}
