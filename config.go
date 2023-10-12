package loggingdrain

type minerConfig struct {
	Mask  maskConfig
	Drain drainConfig
}

type drainConfig struct {
	Similarity  float32
	Depth       int
	MaxChildren int
	MaxCluster  int
}

type maskConfig struct {
	Prefix           string
	Suffix           string
	MaskInstructions []maskInstruction
}

type maskInstruction struct {
	Pattern  string
	MaskWith string
}
