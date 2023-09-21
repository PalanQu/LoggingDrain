package config

type MaskConfig struct {
	Prefix           string            `mapstructure:"prefix" defaultvalue:"[:"`
	Suffix           string            `mapstructure:"suffix" defaultvalue:":]"`
	MaskInstructions []MaskInstruction `mapstructure:"instructions"`
}

type MaskInstruction struct {
	Pattern  string `mapstructure:"pattern"`
	MaskWith string `mapstructure:"mask_with"`
}
