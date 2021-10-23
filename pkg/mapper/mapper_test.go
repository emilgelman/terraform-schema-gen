package mapper

import (
	"testing"

	"github.com/go-openapi/jsonreference"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"k8s.io/kube-openapi/pkg/common"
	"k8s.io/kube-openapi/pkg/validation/spec"
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

func TestMapArrayOfPrimitive(t *testing.T) {
	m := New()
	definitions := map[string]common.OpenAPIDefinition{
		"car": {
			Schema: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{"object"},
					Properties: map[string]spec.Schema{
						"colors": {
							SchemaProps: spec.SchemaProps{
								Type: []string{"array"},
								Items: &spec.SchemaOrArray{
									Schema: &spec.Schema{
										SchemaProps: spec.SchemaProps{
											Default: "",
											Type:    []string{"string"},
											Format:  "",
										},
									},
								},
							},
						},
					},
					Required: []string{"colors"},
				},
			},
		},
	}
	output := m.Map(definitions)
	expected := map[string]map[string]*schema.Schema{
		"car": {
			"colors": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	assert.Equal(t, expected, output)

}

func TestMapArrayOfObject(t *testing.T) {
	m := New()
	ref := func(path string) spec.Ref {
		return spec.Ref{Ref: jsonreference.MustCreateRef(path)}
	}
	definitions := map[string]common.OpenAPIDefinition{
		"engineSpec": {
			Schema: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{"object"},
					Properties: map[string]spec.Schema{
						"cylinders": {
							SchemaProps: spec.SchemaProps{
								Type: []string{"array"},
								Items: &spec.SchemaOrArray{
									Schema: &spec.Schema{
										SchemaProps: spec.SchemaProps{
											Default: map[string]interface{}{},
											Ref:     ref("cylinder"),
										},
									},
								},
							},
						},
					},
					Required: []string{"cylinders"},
				},
			},
			Dependencies: []string{
				"cylinder"},
		},
		"cylinder": {
			Schema: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{"object"},
					Properties: map[string]spec.Schema{
						"Number": {
							SchemaProps: spec.SchemaProps{
								Default: "",
								Type:    []string{"string"},
								Format:  "",
							},
						},
					},
					Required: []string{"Number"},
				},
			},
		},
	}
	output := m.Map(definitions)
	expected := map[string]map[string]*schema.Schema{
		"engineSpec": {
			"cylinders": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &map[string]*schema.Schema{
					"Number": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"cylinder": {
			"Number": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			}},
	}
	assert.Equal(t, expected, output)
}
