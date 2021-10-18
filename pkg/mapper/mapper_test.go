package mapper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"k8s.io/kube-openapi/pkg/common"
	"k8s.io/kube-openapi/pkg/validation/spec"
	"testing"
)

func TestMapper(t *testing.T) {
	m := New()
	definitions := map[string]common.OpenAPIDefinition{
		"car": {
			Schema: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{"object"},
					Properties: map[string]spec.Schema{
						"model": {
							SchemaProps: spec.SchemaProps{
								Description: "Describes the car model",
								Type:        []string{"string"},
							},
						},
					},
					Required: []string{"model"},
				},
			},
		},
	}
	output := m.Map(definitions)
	expected := map[string]map[string]*schema.Schema{
		"car": {
			"model": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Describes the car model",
			},
		},
	}
	assert.Equal(t, expected, output)
}
