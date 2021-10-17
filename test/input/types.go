package input

// +k8s:openapi-gen=true
type Car struct {
	Make       string
	Model      string
	EngineSpec EngineSpec
}

// +k8s:openapi-gen=true
type EngineSpec struct {
	BHP string
}
