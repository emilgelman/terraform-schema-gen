package cmd

import (
	"fmt"
	"github.com/emilgelman/terraform-schema-gen/pkg/generator"
	"github.com/spf13/cobra"
)

var (
	input, output string
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "gen",
	Long:  `gen`,
	Run:   nil,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(input, output)
		g := generator.New(input, output)
		if err := g.Parse(); err != nil {
			return err
		}
		if err := g.Export(); err != nil {
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
	rootCmd.AddCommand(genCmd)
}
