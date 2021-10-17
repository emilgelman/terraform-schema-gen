package generator

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hexops/valast"
)

func (g *Generator) Export() error {
	entries, err := g.createTerraformSchemaEntries()
	if err != nil {
		return err
	}
	t := template.Must(template.New("tfSchemasTemplate").Parse(tfSchemasTemplate))
	var buffer bytes.Buffer
	final := tfSchemas{Schemas: strings.Join(entries, "\n"), Package: g.outputPackage}
	if err := t.Execute(&buffer, final); err != nil {
		panic(err)
	}
	file, err := ioutil.ReadAll(&buffer)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(g.output, file, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) createTerraformSchemaEntries() ([]string, error) {
	entries := make([]string, 0, len(g.schemas))
	for k, v := range g.schemas {
		name := g.formatName(k)
		s := g.formatSchema(v)
		tfs := tfSchema{Name: name, Params: s}
		var buffer bytes.Buffer
		if err := template.Must(template.New("schemaTemplate").Parse(schemaTemplate)).Execute(&buffer, tfs); err != nil {
			panic(err)
		}
		entries = append(entries, buffer.String())
	}
	return entries, nil
}

func (g *Generator) formatName(name string) string {
	tmp := strings.Split(name, ".")
	return tmp[len(tmp)-1]
}

func (g *Generator) formatSchema(schema map[string]*schema.Schema) string {
	s := valast.String(schema)
	return fixValueTypeEnum(s)
}

func fixValueTypeEnum(params string) string {
	return tfValueTypeRegex.ReplaceAllString(params, "schema.ValueType(schema.$1)")
}
