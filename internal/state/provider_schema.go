package state

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-registry-address"
)

type ProviderSchema struct {
	Address tfaddr.Provider
	Version *version.Version
	Source  SchemaSource

	Schema *Schema
}

func (ps *ProviderSchema) idAsArgs() []interface{} {
	args := make([]interface{}, 0)

	args = append(args, ps.Address)
	args = append(args, ps.Version)
	args = append(args, ps.Source)

	return args
}

type ProviderSchemaIdIdx struct{}

func composeProviderSchemaIndex(addr tfaddr.Provider, pv *version.Version, src SchemaSource) string {
	return fmt.Sprintf("%s@%s@%s",
		addr.String(),
		pv.String(),
		src.IndexValue())
}

func (i *ProviderSchemaIdIdx) FromObject(obj interface{}) (bool, []byte, error) {
	ps, ok := obj.(*ProviderSchema)
	if !ok {
		return false, nil, fmt.Errorf("object is not ProviderSchema")
	}

	idx := composeProviderSchemaIndex(ps.Address, ps.Version, ps.Source)
	return true, []byte(idx), nil
}

func (i *ProviderSchemaIdIdx) FromArgs(args ...interface{}) ([]byte, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("expected exactly 3 args, %d given", len(args))
	}

	addr, ok := args[0].(tfaddr.Provider)
	if !ok {
		return nil, fmt.Errorf("unexpected argument type for address")
	}
	pv, ok := args[1].(*version.Version)
	if !ok {
		return nil, fmt.Errorf("unexpected argument type for version")
	}
	source, ok := args[2].(SchemaSource)
	if !ok {
		return nil, fmt.Errorf("unexpected argument type for source")
	}

	idx := composeProviderSchemaIndex(addr, pv, source)
	return []byte(idx), nil
}

func (i *ProviderSchemaIdIdx) PrefixFromArgs(args ...interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expected exactly 1 arg, %d given", len(args))
	}

	addr, ok := args[0].(tfaddr.Provider)
	if !ok {
		return nil, fmt.Errorf("unexpected argument type for address")
	}

	prefix := fmt.Sprintf("%s@", addr.String())
	return []byte(prefix), nil
}

const providerSchemaTableName = "provider_schema"

var providerSchemaTable = &memdb.TableSchema{
	Name: providerSchemaTableName,
	Indexes: map[string]*memdb.IndexSchema{
		"id": {
			Name:    "id",
			Unique:  true,
			Indexer: &ProviderSchemaIdIdx{},
		},
	},
}

func (s *StateStore) AddProviderSchema(ps *ProviderSchema) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	// TODO: Introduce Exists method to Txn?
	obj, err := txn.First(providerSchemaTableName, "id", ps.idAsArgs()...)
	if err != nil {
		return err
	}
	if obj != nil {
		return &AlreadyExistsError{
			Idx: composeProviderSchemaIndex(ps.Address, ps.Version, ps.Source),
		}
	}

	err = txn.Insert(providerSchemaTableName, ps)
	if err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s *StateStore) SchemaFor(modPath string, addr tfaddr.Provider, vc version.Constraints) (*Schema, error) {
	txn := s.db.Txn(false)

	_, err := txn.Get(providerSchemaTableName, "id_prefix", addr)
	if err != nil {
		return nil, err
	}

	// TODO: Pick schema based on:
	// - version constraints
	// - source (locally sourced schema is preferred;
	// 		when two locally sourced schemas of different version are available,
	// 		hierarchy proximity is considered, as also explained in #180)
	// - higher versions are generally preferred

	// it.Next()

	return nil, nil
}
