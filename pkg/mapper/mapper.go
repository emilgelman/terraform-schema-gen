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

func (m *Mapper) Convert(definitions map[string]common.OpenAPIDefinition) map[string]map[string]*schema.Schema {
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
		s := make(map[string]*schema.Schema)
		m.parseDefinition(definition.Name, definition.Name, &definition.Definition.Schema, s)
		schemas[definition.Name] = s
	}
	return schemas
}

func (m *Mapper) parseDefinition(rootName, name string, s *spec.Schema, wtf map[string]*schema.Schema) {
	for n := range s.Properties {
		prop := s.Properties[n]
		if prop.SchemaProps.Type == nil {
			path := prop.Ref.Ref.GetURL().Path
			ss := m.schemas[path]
			wtf[n] = &schema.Schema{Type: schema.TypeList, Elem: &schema.Resource{Schema: ss}}
			continue
		}
		m.parseDefinition(rootName, n, &prop, m)
	}
	if name == rootName {
		return
	}
	if s.Type == nil {
		return
	}
	tType := s.Type[0]
	switch tType {
	case "object":
		wtf[name] = &schema.Schema{Type: schema.TypeMap}
	case "string":
		wtf[name] = &schema.Schema{Type: schema.TypeString}
	}
}
