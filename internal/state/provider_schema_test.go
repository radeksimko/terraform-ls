package state

import (
	"errors"
	"testing"

	"github.com/hashicorp/go-version"
	tfaddr "github.com/hashicorp/terraform-registry-address"
)

func TestProviderSchema_validate(t *testing.T) {
	err := providerSchemaTable.Validate()
	if err != nil {
		t.Fatal(err)
	}
}

func TestProviderSchema_duplicateEntries(t *testing.T) {
	s, err := NewStateStore()
	if err != nil {
		t.Fatal(err)
	}

	ps := &ProviderSchema{
		Address: tfaddr.Provider{
			Hostname:  tfaddr.DefaultRegistryHost,
			Namespace: "hashicorp",
			Type:      "aws",
		},
		Version: version.Must(version.NewVersion("1.0.0")),
		Source:  PreloadedSchemaSource{},
		Schema:  &Schema{},
	}

	err = s.AddProviderSchema(ps)
	if err != nil {
		t.Fatal(err)
	}

	err = s.AddProviderSchema(ps)
	if err == nil {
		t.Fatal("expected duplicate insertion to fail")
	}

	aeErr := &AlreadyExistsError{}
	if !errors.As(err, &aeErr) {
		t.Fatalf("unexpected error: %s", err)
	}
}
