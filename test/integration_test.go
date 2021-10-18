package test

import (
	"os/exec"
	"plugin"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

var expectedSchemas = map[string]map[string]*schema.Schema{
	"GetEngineSpecSchema": {"bhp": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}},
	"GetCarSchema": {
		"engineSpec": &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"bhp": {Type: schema.TypeString, Required: true},
			}},
		},
		"make":  &schema.Schema{Type: schema.TypeString, Required: true},
		"model": &schema.Schema{Type: schema.TypeString},
	},
}

func TestGenerator(t *testing.T) {
	generateOpenAPISchema(t)
	compileOpenAPISchemaAsPlugin(t)
	buildTerraformSchemaGen(t)
	runGenerator(t)
	compileGeneratedAsPlugin(t)

	plug, err := plugin.Open("./output/main/terraform_generated.so")
	assert.NoError(t, err)
	for name, expected := range expectedSchemas {
		schemaFunc, err := plug.Lookup(name)
		assert.NoError(t, err)
		schemaSupplier := schemaFunc.(func() map[string]*schema.Schema)
		assert.Equal(t, expected, schemaSupplier())
	}
}

func compileGeneratedAsPlugin(t *testing.T) {
	command := exec.Command("go", "build", "-buildmode=plugin",
		"-o", "./output/main/terraform_generated.so", "./output/main/terraform_generated.go")
	err := command.Run()
	assert.NoError(t, err)
}

func runGenerator(t *testing.T) {
	command := exec.Command("./output/terraform-schema-gen", "gen",
		"--input", "./output/main/openapi_generated.so", "--output", "./output/main/terraform_generated.go", "--package", "main")
	err := command.Run()
	assert.NoError(t, err)
}

func buildTerraformSchemaGen(t *testing.T) {
	command := exec.Command("go", "build", "-o", "./output/terraform-schema-gen", "../main.go")
	err := command.Run()
	assert.NoError(t, err)
}

func compileOpenAPISchemaAsPlugin(t *testing.T) {
	command := exec.Command("go", "build", "-buildmode=plugin",
		"-o", "./output/main/openapi_generated.so", "./output/main/openapi_generated.go")
	err := command.Run()
	assert.NoError(t, err)
}

func generateOpenAPISchema(t *testing.T) {
	command := exec.Command("go", "run", "k8s.io/kube-openapi/cmd/openapi-gen",
		"-i", "./input", "-p", "./output/main", "-h", "./header.txt", "-o", ".")
	err := command.Run()
	assert.NoError(t, err)
}
