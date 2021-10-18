package generator

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"k8s.io/kube-openapi/pkg/common"
)

type OpenAPIDefinitionLoader interface {
	Load() (map[string]common.OpenAPIDefinition, error)
}
type OpenAPIDefinitionMapper interface {
	Map(map[string]common.OpenAPIDefinition) map[string]map[string]*schema.Schema
}

type SchemaExporter interface {
	Export(map[string]map[string]*schema.Schema) error
}

type Generator struct {
	definitionLoader OpenAPIDefinitionLoader
	mapper           OpenAPIDefinitionMapper
	exporter         SchemaExporter
}

func New(definitionLoader OpenAPIDefinitionLoader, mapper OpenAPIDefinitionMapper, exporter SchemaExporter) *Generator {
	return &Generator{definitionLoader: definitionLoader, mapper: mapper, exporter: exporter}
}

func (g *Generator) Generate() error {
	definitions, err := g.definitionLoader.Load()
	if err != nil {
		return err
	}
	schemas := g.mapper.Map(definitions)
	if err = g.exporter.Export(schemas); err != nil {
		return err
	}
	return nil
}
