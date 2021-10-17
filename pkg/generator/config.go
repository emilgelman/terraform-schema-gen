package generator

type Config struct {
	input         string
	output        string
	outputPackage string
}

func NewConfig(input, output, outputPackage string) Config {
	return Config{
		input:         input,
		output:        output,
		outputPackage: outputPackage,
	}
}
