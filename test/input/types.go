package input

// +k8s:openapi-gen=true
type Car struct {
	Make       string     `json:"make"`
	Model      string     `json:"model,omitempty"`
	EngineSpec EngineSpec `json:"engineSpec"`
}

// +k8s:openapi-gen=true
type EngineSpec struct {
	BHP string `json:"bhp"`
}
