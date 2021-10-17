package exporter

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hexops/valast"
)

var tfValueTypeRegex = regexp.MustCompile(`schema.ValueType\((.*)\)`)

type Exporter struct {
	output        string
	outputPackage string
}

func (e *Exporter) Export(schemas map[string]map[string]*schema.Schema) error {
	entries, err := e.createTerraformSchemaEntries(schemas)
	if err != nil {
		return err
	}
	t := template.Must(template.New("tfSchemasTemplate").Parse(tfSchemasTemplate))
	var buffer bytes.Buffer
	final := tfSchemas{Schemas: strings.Join(entries, "\n"), Package: e.outputPackage}
	if err := t.Execute(&buffer, final); err != nil {
		panic(err)
	}
	file, err := ioutil.ReadAll(&buffer)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(e.output, file, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (e *Exporter) createTerraformSchemaEntries(schemas map[string]map[string]*schema.Schema) ([]string, error) {
	entries := make([]string, 0, len(schemas))
	for k, v := range schemas {
		name := e.formatName(k)
		s := e.formatSchema(v)
		tfs := tfSchema{Name: name, Params: s}
		var buffer bytes.Buffer
		if err := template.Must(template.New("schemaTemplate").Parse(schemaTemplate)).Execute(&buffer, tfs); err != nil {
			panic(err)
		}
		entries = append(entries, buffer.String())
	}
	return entries, nil
}

func (e *Exporter) formatName(name string) string {
	tmp := strings.Split(name, ".")
	return tmp[len(tmp)-1]
}

func (e *Exporter) formatSchema(schema map[string]*schema.Schema) string {
	s := valast.String(schema)
	return fixValueTypeEnum(s)
}

func fixValueTypeEnum(params string) string {
	return tfValueTypeRegex.ReplaceAllString(params, "schema.ValueType(schema.$1)")
}
