# terraform-schema-gen

This repository contains a CLI to
generate [Terraform schema](https://www.terraform.io/docs/extend/schemas/schema-types.html) out of Go structs.

The generator will convert nested Go structs (at any level) to the equivalent Terraform schema.

The schemas' `Required` or `Optional` property is set based on the json `omitempty` tag.
If omitempty is set, the property is marked as Optional, else as Required.

### Known limitations

Terraform functions and validations are not supported, as there is no current way to express them from struct properties.

## Usage

1. Run the terraform-schema-gen CLI to convert a directory containing Go structs to a Terraform schema:
```
go run github.com/emilgelman/terraform-schema-gen gen --input <input directory> --output terraform_schema_generated.go --package <output package name>
```

The command can be bundled in a `go generate` for automation, for example:
```go
package generate

//go:generate go run github.com/emilgelman/terraform-schema-gen gen --input ./input --output ./output/v1/terraform_schema_generated.go --package v1
```

