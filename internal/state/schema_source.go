package state

import "fmt"

type schemaSrcSigil struct{}

type SchemaSource interface {
	isSchemaSrcImpl() schemaSrcSigil
	IndexValue() string
}

type PreloadedSchemaSource struct {
}

func (PreloadedSchemaSource) isSchemaSrcImpl() schemaSrcSigil {
	return schemaSrcSigil{}
}

func (PreloadedSchemaSource) IndexValue() string {
	return "preloaded"
}

type LocalSchemaSource struct {
	ModulePath string
}

func (LocalSchemaSource) isSchemaSrcImpl() schemaSrcSigil {
	return schemaSrcSigil{}
}

func (lss LocalSchemaSource) IndexValue() string {
	return fmt.Sprintf("local(%s)", lss.ModulePath)
}
