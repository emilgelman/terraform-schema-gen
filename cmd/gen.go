package cmd

import (
	"github.com/spf13/cobra"

	"github.com/emilgelman/terraform-schema-gen/pkg/generator"
)

var (
	input, output, outputPackage string
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "gen",
	Long:  `gen`,
	Run:   nil,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := generator.NewConfig(input, output, outputPackage)
		g := generator.New(nil, config)
		if err := g.Generate(); err != nil {
			return err
		}
		if err := g.Export(nil); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	genCmd.Flags().StringVarP(&input, "input", "i", "", "input file")
	genCmd.MarkFlagRequired("input")
	genCmd.Flags().StringVarP(&output, "output", "o", "", "output file")
	genCmd.MarkFlagRequired("output")
	genCmd.Flags().StringVarP(&outputPackage, "package", "p", "", "output package")
	genCmd.MarkFlagRequired("package")
	rootCmd.AddCommand(genCmd)
}
