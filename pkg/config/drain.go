package config

type DrainConfig struct {
	Similarity  float32 `mapstructure:"similarity" defaultvalue:"0.4"`
	Depth       int     `mapstructure:"depth" defaultvalue:"5"`
	MaxChildren int     `mapstructure:"max_children" defaultvalue:"100"`
	MaxCluster  int     `mapstructure:"max_children" defaultvalue:"1000"`
}
