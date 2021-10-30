package generator

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"k8s.io/gengo/types"
)

type StructParser interface {
	Parse() ([]*types.Type, error)
}
type OpenAPIDefinitionMapper interface {
	Map([]*types.Type) map[string]map[string]*schema.Schema
}

type SchemaExporter interface {
	Export(map[string]map[string]*schema.Schema) error
}

type Generator struct {
	parser   StructParser
	mapper   OpenAPIDefinitionMapper
	exporter SchemaExporter
}

func New(parser StructParser, mapper OpenAPIDefinitionMapper, exporter SchemaExporter) *Generator {
	return &Generator{parser: parser, mapper: mapper, exporter: exporter}
}

func (g *Generator) Generate() error {
	parsedTypes, err := g.parser.Parse()
	if err != nil {
		return err
	}
	schemas := g.mapper.Map(parsedTypes)
	if err := g.exporter.Export(schemas); err != nil {
		return err
	}
	return nil
}
