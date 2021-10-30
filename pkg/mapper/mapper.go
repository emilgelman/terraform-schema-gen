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
		name := parsedTypes[t].Name.Name
		traverse(parsedTypes[t], name, name, "", s)
		res[name] = s
	}
	return res
}

func traverse(t *types.Type, name, rootName, tags string, s map[string]*schema.Schema) {
	name = strings.ToLower(name)
	rootName = strings.ToLower(rootName)
	switch t.Kind {
	// The first cases handles nested structures and translates them recursively
	// If it is a pointer we need to unwrap and recurse
	case types.Pointer:
		traverse(t.Elem, name, rootName, "", s)

	// If it is a struct we translate each field
	case types.Struct:
		x := make(map[string]*schema.Schema)
		for _, member := range t.Members {
			traverse(member.Type, member.Name, rootName, member.Tags, x)
		}
		if name == rootName {
			for k, v := range x {
				s[k] = v
			}
			return
		}
		converted := &schema.Schema{
			Type: schema.TypeList,
			Elem: &schema.Resource{Schema: x},
		}
		setOptionalOrRequired(converted, tags)
		s[name] = converted

	// If it is a slice we translate the inner element type
	case types.Slice:
		x := make(map[string]*schema.Schema)
		traverse(t.Elem, t.Name.Name, rootName, "", x)

		if t.Elem.Kind == types.Builtin {
			converted := &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{Type: schema.TypeString}, //TODO: handle other types of primitive slices
			}
			setOptionalOrRequired(converted, tags)
			s[name] = converted
			return
		}
		s[name] = x[strings.ToLower(t.Name.Name)]

	// If it is a map return map[string]string //TODO: handle complex map structures
	case types.Map:
		converted := &schema.Schema{
			Type: schema.TypeMap,
			Elem: &schema.Schema{Type: schema.TypeString},
		}

		setOptionalOrRequired(converted, tags)
		s[name] = converted

	// Otherwise we cannot traverse anywhere so this finishes the the recursion
	// If it is a builtin type translate it
	case types.Builtin:
		converted := convertBuiltinType(t, tags)
		s[name] = converted

	default:
		fmt.Printf("Unknown type %+v", t)
	}
}

func convertBuiltinType(t *types.Type, tags string) *schema.Schema {
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
	converted := &schema.Schema{
		Type: schemaType,
	}
	setOptionalOrRequired(converted, tags)
	return converted
}

func setOptionalOrRequired(s *schema.Schema, tags string) {
	if strings.Contains(tags, "omitempty") {
		s.Optional = true
	} else {
		s.Required = true
	}
}
