package mapper

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"k8s.io/gengo/types"
)

type Mapper struct {
}

func New() *Mapper {
	return &Mapper{}
}

func (m *Mapper) Map(parsedTypes []*types.Type) map[string]map[string]*schema.Schema {
	res := make(map[string]map[string]*schema.Schema)
	for t := range parsedTypes {
		s := make(map[string]*schema.Schema)
		traverse(parsedTypes[t], parsedTypes[t].Name.Name, parsedTypes[t].Name.Name, s)
		res[parsedTypes[t].Name.Name] = s
	}
	return res
}

func traverse(t *types.Type, name, rootName string, s map[string]*schema.Schema) {
	name = strings.ToLower(name)
	rootName = strings.ToLower(rootName)
	switch t.Kind {
	// The first cases handles nested structures and translates them recursively

	// If it is a pointer we need to unwrap and recurse
	case types.Pointer:
		traverse(t.Elem, name, rootName, s)

	// If it is a struct we translate each field
	case types.Struct:
		x := make(map[string]*schema.Schema)
		for _, member := range t.Members {
			traverse(member.Type, member.Name, rootName, x)
		}

		if name == rootName {
			for k, v := range x {
				s[k] = v
			}
			return
		}
		s[name] = &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem:     &schema.Resource{Schema: x},
		}

	// If it is a slice we translate the inner element type
	case types.Slice:
		x := make(map[string]*schema.Schema)
		traverse(t.Elem, t.Name.Name, rootName, x)
		if t.Elem.Kind == types.Builtin {
			s[name] = &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString}, //TODO: handle other types of slices
			}
			return
		}
		s[name] = x[strings.ToLower(t.Name.Name)]

	// If it is a map return map[string]string //TODO: handle complex map structures
	case types.Map:
		s[name] = &schema.Schema{
			Type:     schema.TypeMap,
			Required: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		}

	// Otherwise we cannot traverse anywhere so this finishes the the recursion
	// If it is a builtin type translate it
	case types.Builtin:
		var schemaType schema.ValueType
		switch t.Name {
		case types.String.Name:
			schemaType = schema.TypeString
		case types.Int.Name, types.Int32.Name, types.Int64.Name:
			schemaType = schema.TypeInt
		case types.Float.Name, types.Float32.Name, types.Float64.Name:
			schemaType = schema.TypeFloat
		case types.Bool.Name:
			schemaType = schema.TypeBool

		}
		translated := &schema.Schema{
			Type:     schemaType,
			Required: true,
		}
		s[name] = translated

	default:
		fmt.Printf("Unknown type %+v", t)
	}
}
