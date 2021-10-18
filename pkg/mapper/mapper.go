package mapper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"k8s.io/kube-openapi/pkg/common"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

type Mapper struct {
}

type SchemaDefinition struct {
	Name       string
	Definition common.OpenAPIDefinition
}

func New() *Mapper {
	return &Mapper{}
}

func (m *Mapper) Map(definitions map[string]common.OpenAPIDefinition) map[string]map[string]*schema.Schema {
	stack := m.createDefinitionsStack(definitions)
	return m.parseDefinitionsStack(stack)
}

func (m *Mapper) createDefinitionsStack(definitions map[string]common.OpenAPIDefinition) []SchemaDefinition {
	var stack []SchemaDefinition
	for name, definition := range definitions {
		stack = append(stack, SchemaDefinition{Name: name, Definition: definition})
		for _, dependency := range definition.Dependencies {
			stack = append(stack, SchemaDefinition{Name: dependency, Definition: definitions[dependency]})
		}
	}
	return stack
}

func (m *Mapper) parseDefinitionsStack(stack []SchemaDefinition) map[string]map[string]*schema.Schema {
	schemas := make(map[string]map[string]*schema.Schema)
	for len(stack) > 0 {
		definition := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		tfSchema := make(map[string]*schema.Schema)
		m.parseDefinition(definition.Name, definition.Name, &definition.Definition.Schema, tfSchema, schemas)
		schemas[definition.Name] = tfSchema
	}
	return schemas
}

func (m *Mapper) parseDefinition(rootName, name string, openapiSchema *spec.Schema, tfSchema map[string]*schema.Schema, schemas map[string]map[string]*schema.Schema) {
	for i := range openapiSchema.Properties {
		prop := openapiSchema.Properties[i]
		if prop.SchemaProps.Type == nil {
			path := prop.Ref.Ref.GetURL().Path
			ss := schemas[path]
			tfSchema[i] = &schema.Schema{Type: schema.TypeList, Elem: &schema.Resource{Schema: ss}}
			continue
		}
		m.parseDefinition(rootName, i, &prop, tfSchema, schemas)
	}
	if name == rootName {
		for i := range openapiSchema.Required {
			tfSchema[openapiSchema.Required[i]].Required = true
		}
		return
	}
	if openapiSchema.Type == nil {
		return
	}
	newSchema := &schema.Schema{}
	tType := openapiSchema.Type[0]
	switch tType {
	case "object":
		newSchema.Type = schema.TypeMap
	case "string":
		newSchema.Type = schema.TypeString
	}
	newSchema.Description = openapiSchema.Description
	tfSchema[name] = newSchema
}
