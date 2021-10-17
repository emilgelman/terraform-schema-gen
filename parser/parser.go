package parser

import (
	"bytes"
	"io/ioutil"
	"k8s.io/kube-openapi/pkg/common"
	"plugin"
	"regexp"
	"strings"
	"text/template"

	"github.com/go-openapi/jsonreference"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hexops/valast"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

type Parser struct {
	input   string
	output  string
	schemas map[string]map[string]*schema.Schema
}

var tfValueTypeRegex = regexp.MustCompile(`schema.ValueType\((.*)\)`)

func New(input, output string) *Parser {
	return &Parser{schemas: make(map[string]map[string]*schema.Schema), input: input, output: output}
}
func (p *Parser) Parse() error {
	plug, err := plugin.Open(p.input)
	if err != nil {
		return err
	}
	lookup, err := plug.Lookup("GetOpenAPIDefinitions")
	if err != nil {
		return err
	}
	openapiDefinitionsSupplierFunc := lookup.(func(common.ReferenceCallback) map[string]common.OpenAPIDefinition)

	definitions := openapiDefinitionsSupplierFunc(func(path string) spec.Ref {
		return spec.Ref{Ref: jsonreference.MustCreateRef(path)}
	})
	for name := range definitions {
		definition := definitions[name]
		for _, dependency := range definition.Dependencies {
			m := make(map[string]*schema.Schema)
			s := definitions[dependency].Schema
			p.parseDefinition(dependency, dependency, &s, m)
			p.schemas[dependency] = m
		}
		m := make(map[string]*schema.Schema)
		p.parseDefinition(name, name, &definition.Schema, m)
		p.schemas[name] = m
	}
	return nil
}

func (p *Parser) Export() {
	entries := make([]string, 0, len(p.schemas))
	for k, v := range p.schemas {
		arr := strings.Split(k, ".")
		name := arr[len(arr)-1]
		params := valast.String(v)
		params = fixValueTypeEnum(params)
		rt := tfSchema{Name: name, Params: params}

		var result bytes.Buffer
		if err := template.Must(template.New("tfValueTypeRegex").Parse(schemaTemplate)).Execute(&result, rt); err != nil {
			panic(err)
		}
		str := result.String()
		entries = append(entries, str)
	}

	t := template.Must(template.New("validate").Parse(tfSchemasTemplate))
	var result bytes.Buffer
	final := tfSchemas{Schemas: strings.Join(entries, "\n")}
	if err := t.Execute(&result, final); err != nil {
		panic(err)
	}
	file, err := ioutil.ReadAll(&result)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(p.output, file, 0600)
	if err != nil {
		panic(err)
	}
}

func fixValueTypeEnum(params string) string {
	return tfValueTypeRegex.ReplaceAllString(params, "schema.ValueType(schema.$1)")
}

func (p *Parser) parseDefinition(rootName, name string, s *spec.Schema, m map[string]*schema.Schema) {
	for n := range s.Properties {
		prop := s.Properties[n]
		if prop.SchemaProps.Type == nil {
			path := prop.Ref.Ref.GetURL().Path
			ss := p.schemas[path]
			m[n] = &schema.Schema{Type: schema.TypeList, Elem: &schema.Resource{Schema: ss}}
			continue
		}
		p.parseDefinition(rootName, n, &prop, m)
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
