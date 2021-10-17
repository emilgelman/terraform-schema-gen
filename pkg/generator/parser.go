package generator

import (
	"plugin"
	"regexp"

	"github.com/go-openapi/jsonreference"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"k8s.io/kube-openapi/pkg/common"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

type Generator struct {
	input         string
	output        string
	outputPackage string
	schemas       map[string]map[string]*schema.Schema
}

type SchemaDefinition struct {
	Name       string
	Definition common.OpenAPIDefinition
}

var tfValueTypeRegex = regexp.MustCompile(`schema.ValueType\((.*)\)`)

func New(input, output, outputPackage string) *Generator {
	return &Generator{schemas: make(map[string]map[string]*schema.Schema), input: input, output: output, outputPackage: outputPackage}
}
func (g *Generator) Parse() error {
	definitions, err := g.getDefinitions()
	if err != nil {
		return err
	}
	stack := g.createDefinitionsStack(definitions)
	g.parseDefinitionsStack(stack)
	return nil
}

func (g *Generator) getDefinitions() (map[string]common.OpenAPIDefinition, error) {
	plug, err := plugin.Open(g.input)
	if err != nil {
		return nil, err
	}
	lookup, err := plug.Lookup("GetOpenAPIDefinitions")
	if err != nil {
		return nil, err
	}
	openapiDefinitionsSupplierFunc := lookup.(func(common.ReferenceCallback) map[string]common.OpenAPIDefinition)
	definitions := openapiDefinitionsSupplierFunc(func(path string) spec.Ref {
		return spec.Ref{Ref: jsonreference.MustCreateRef(path)}
	})
	return definitions, nil
}

func (g *Generator) parseDefinitionsStack(stack []SchemaDefinition) {
	for len(stack) > 0 {
		definition := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		m := make(map[string]*schema.Schema)
		g.parseDefinition(definition.Name, definition.Name, &definition.Definition.Schema, m)
		g.schemas[definition.Name] = m
	}
}

func (g *Generator) createDefinitionsStack(definitions map[string]common.OpenAPIDefinition) []SchemaDefinition {
	var stack []SchemaDefinition
	for name, definition := range definitions {
		stack = append(stack, SchemaDefinition{Name: name, Definition: definition})
		for _, dependency := range definition.Dependencies {
			stack = append(stack, SchemaDefinition{Name: dependency, Definition: definitions[dependency]})
		}
	}
	return stack
}

func (g *Generator) parseDefinition(rootName, name string, s *spec.Schema, m map[string]*schema.Schema) {
	for n := range s.Properties {
		prop := s.Properties[n]
		if prop.SchemaProps.Type == nil {
			path := prop.Ref.Ref.GetURL().Path
			ss := g.schemas[path]
			m[n] = &schema.Schema{Type: schema.TypeList, Elem: &schema.Resource{Schema: ss}}
			continue
		}
		g.parseDefinition(rootName, n, &prop, m)
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
		m[name] = &schema.Schema{Type: schema.TypeMap}
	case "string":
		m[name] = &schema.Schema{Type: schema.TypeString}
	}
}
