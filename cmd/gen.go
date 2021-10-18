package cmd

import (
	"github.com/emilgelman/terraform-schema-gen/pkg/exporter"
	"github.com/emilgelman/terraform-schema-gen/pkg/mapper"
	"github.com/emilgelman/terraform-schema-gen/pkg/openapi"
	"github.com/spf13/cobra"

	"github.com/emilgelman/terraform-schema-gen/pkg/generator"
)

var config generator.Config

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "gen",
	Long:  `gen`,
	Run:   nil,
	RunE: func(cmd *cobra.Command, args []string) error {
		loader := openapi.New(config.Input)
		mapper := mapper.New()
		exporter := exporter.New(config.Output, config.OutputPackage)
		g := generator.New(loader, mapper, exporter)
		return g.Generate()
	},
}

func init() {
	genCmd.Flags().StringVarP(&config.Input, "input", "i", "", "input file")
	genCmd.MarkFlagRequired("input")
	genCmd.Flags().StringVarP(&config.Output, "output", "o", "", "output file")
	genCmd.MarkFlagRequired("output")
	genCmd.Flags().StringVarP(&config.OutputPackage, "package", "p", "", "output package")
	genCmd.MarkFlagRequired("package")
	rootCmd.AddCommand(genCmd)
}
