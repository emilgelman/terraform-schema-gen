package openapi

import (
	"plugin"

	"github.com/go-openapi/jsonreference"
	"k8s.io/kube-openapi/pkg/common"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

type Loader struct {
	input string
}

func New(input string) *Loader {
	return &Loader{input: input}
}

func (l *Loader) Load() (map[string]common.OpenAPIDefinition, error) {
	plug, err := plugin.Open(l.input)
	if err != nil {
		return nil, err
	}
	lookup, err := plug.Lookup("GetOpenAPIDefinitions")
	if err != nil {
		return nil, err
	}
	openapiDefinitionsSupplierFunc := lookup.(func(common.ReferenceCallback) map[string]common.OpenAPIDefinition)
	definitions := openapiDefinitionsSupplierFunc(func(path string) spec.Ref {
		return spec.Ref{Ref: jsonreference.MustCreateRef(path)}
	})
	return definitions, nil
}
