package mapper

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"k8s.io/gengo/types"
)

// nolint
func TestMapper(t *testing.T) {
	var tests = []struct {
		name     string
		input    []*types.Type
		expected map[string]map[string]*schema.Schema
	}{
		{
			name: "struct with primitive properties",
			input: []*types.Type{
				{
					Name: types.Name{Name: "Car"},
					Kind: types.Struct,
					Members: []types.Member{
						{
							Name: "Model",
							Type: types.String,
						},
						{
							Name: "Year",
							Type: types.Int64,
						},
						{
							Name: "IsNew",
							Type: types.Bool,
						},
					},
				},
			},
			expected: map[string]map[string]*schema.Schema{
				"Car": {
					"model": {
						Type:     schema.TypeString,
						Required: true,
					},
					"year": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"isnew": {
						Type:     schema.TypeBool,
						Required: true,
					},
				},
			},
		},
		{
			name: "nested structs",
			input: []*types.Type{
				{
					Name: types.Name{Name: "Car"},
					Kind: types.Struct,
					Members: []types.Member{
						{
							Name: "EngineSpec",
							Type: &types.Type{
								Name: types.Name{Name: "EngineSpec"},
								Kind: types.Struct,
								Members: []types.Member{
									{
										Name: "BHP",
										Type: types.String,
									},
								},
							},
						},
					},
				},
			},
			expected: map[string]map[string]*schema.Schema{
				"Car": {
					"enginespec": {
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"bhp": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "struct with struct slice",
			input: []*types.Type{
				{
					Name: types.Name{Name: "Car"},
					Kind: types.Struct,
					Members: []types.Member{
						{
							Name: "EngineSpec",
							Type: &types.Type{
								Name: types.Name{Name: "EngineSpec"},
								Kind: types.Struct,
								Members: []types.Member{
									{
										Name: "BHP",
										Type: types.String,
									},
								},
							},
						},
					},
				},
			},
			expected: map[string]map[string]*schema.Schema{
				"Car": {
					"enginespec": {
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"bhp": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "struct with primitive slice",
			input: []*types.Type{
				{
					Name: types.Name{Name: "Car"},
					Kind: types.Struct,
					Members: []types.Member{
						{
							Name: "EngineSpec",
							Type: &types.Type{
								Name: types.Name{Name: "[]string"},
								Kind: types.Slice,
								Elem: &types.Type{
									Name: types.Name{Name: "string"},
									Kind: types.Builtin,
								},
							},
						},
					},
				},
			},
			expected: map[string]map[string]*schema.Schema{
				"Car": {
					"enginespec": {
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		{
			name: "struct with map[string]interface{}",
			input: []*types.Type{
				{
					Name: types.Name{Name: "Car"},
					Kind: types.Struct,
					Members: []types.Member{
						{
							Name: "EngineSpec",
							Type: &types.Type{
								Name: types.Name{Name: "map[string]interface{}"},
								Kind: types.Map,
								Elem: &types.Type{
									Name: types.Name{Name: "interface{}"},
									Kind: types.Interface,
								},
							},
						},
					},
				},
			},
			expected: map[string]map[string]*schema.Schema{
				"Car": {
					"enginespec": {
						Type:     schema.TypeMap,
						Required: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
	}
	m := New()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := m.Map(test.input)
			assert.Equal(t, test.expected, output)
		})
	}
}
