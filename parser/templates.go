package parser

type tfSchema struct {
	Name   string
	Params string
}

type tfSchemas struct {
	Schemas string
}

var tfSchemasTemplate = `// +build !ignore_autogenerated

// Code generated by catalog-specification. DO NOT EDIT.

// This file was autogenerated by catalog-specification. Do not edit it manually!

package v1

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

{{.Schemas}}
`

var schemaTemplate = `func Get{{.Name}}Schema() map[string]*schema.Schema {
return {{.Params}}
}
`
