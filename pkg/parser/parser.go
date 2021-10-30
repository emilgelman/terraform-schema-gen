package parser

import (
	"k8s.io/gengo/namer"
	gengo "k8s.io/gengo/parser"
	"k8s.io/gengo/types"
)

type Parser struct {
	input string
}

func New(input string) *Parser {
	return &Parser{input: input}
}

func (p *Parser) Parse() ([]*types.Type, error) {
	parser := gengo.New()
	if err := parser.AddDirRecursive(p.input); err != nil {
		return nil, err
	}
	parsedTypes, err := parser.FindTypes()
	if err != nil {
		return nil, err
	}
	orderer := namer.Orderer{Namer: namer.NewPublicNamer(1)}
	o := orderer.OrderUniverse(parsedTypes)
	return filterStructs(o), nil
}

func filterStructs(typesSlice []*types.Type) []*types.Type {
	var res []*types.Type
	for i := range typesSlice {
		if typesSlice[i].Kind == types.Struct {
			res = append(res, typesSlice[i])
		}
	}
	return res
}
