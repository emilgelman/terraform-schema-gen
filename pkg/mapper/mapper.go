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
	// The first cases handle nested structures and translate them recursively

	// If it is a pointer we need to unwrap and call once again
	case types.Pointer:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := t.Elem
		// Check if the pointer is nil
		// if !originalValue.IsAssignable() {
		//	return
		// }
		// Allocate a new object and set the pointer to it
		// copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		traverse(originalValue, t.Name.Name, rootName, s)
	case types.Interface:
		println("interface")

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
		s[name] = x[strings.ToLower(t.Name.Name)] //TODO: handle primitive slices

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
		if t.Name == types.String.Name {
			translated := &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			}
			s[name] = translated
		}

	default:
		fmt.Println("default")
	}
}
