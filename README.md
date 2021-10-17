# terraform-schema-gen

This repository contains a CLI to
generate [Terraform schema](https://www.terraform.io/docs/extend/schemas/schema-types.html) out of Go structs. The
generator relies on [kube-openapi](https://github.com/kubernetes/kube-openapi) as an intermediate step in the generation
process.

The generator will convert nested structs (at any level) to the equivalent Terraform schema. 

TODO: Required fields

TODO:Description

### Known limitations

Terraform functions and validations are not supported, as this generator relies on Go struct tags.

## Usage

1. Mark your structs with the following comment `// +k8s:openapi-gen=true`
2. Create a header.txt file, to be used with the kube-openapi generator as heading for generated files (can be empty)
2. Run the kube-openapi generator to generate an OpenAPI spec of your structs:

```shell
go run k8s.io/kube-openapi/cmd/openapi-gen -i <input directory> -p <output directory> -h ./header.txt -o <output base>
```

3. Compile the generated OpenAPI spec as a Go plugin

```shell
go build -buildmode=plugin  -o <output_file> <openapi_generated.go file>
```

4. Run the terraform-schema-gen CLI to convert the compiled plugin to Terraform schema:
```
go run github.com/emilgelman/terraform-schema-gen gen --input <openapi_generated.so> --output terraform_schema_generated.go --package <output package name>
```

The entire process can be bundled in a single go file utilizing `go generate`, for example:
```go
package generate

//go:generate go run k8s.io/kube-openapi/cmd/openapi-gen -i ./v1 -p ./output/v1/main -h ./header.txt -o .
//go:generate go build -buildmode=plugin  -o ./output/v1/main/openapi_generated.so ./output/v1/main/openapi_generated.go
//go:generate go run github.com/emilgelman/terraform-schema-gen gen --input ./output/v1/main/openapi_generated.so --output ./output/v1/terraform_schema_generated.go --package v1
```

