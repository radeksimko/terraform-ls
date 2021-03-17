package state

import (
	"github.com/hashicorp/hcl-lang/schema"
)

type Schema struct {
	Provider    *schema.BodySchema
	Resources   map[string]*schema.BodySchema
	DataSources map[string]*schema.BodySchema
}
