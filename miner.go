package loggingdrain

type TemplateMiner struct {
	drain  *drain
	masker *logMasker
}

type LogMessageResponse struct {
	ChangeType    ClusterUpdateType
	Cluster       *LogCluster
	TemplateMined string
	ClusterCount  int
}

func NewTemplateMiner(options ...minerOption) (*TemplateMiner, error) {
	c := newTemplateMinerConfig(options)
	return newTemplateMinerWithConfig(c)
}

func newTemplateMinerWithConfig(config *minerConfig) (*TemplateMiner, error) {
	drain := newDrainWithConfig(config.Drain)
	masker, err := newLogMaskerWithConfig(config.Mask)
	if err != nil {
		return nil, err
	}
	return &TemplateMiner{
		drain:  drain,
		masker: masker,
	}, nil
}

func newTemplateMinerConfig(options []minerOption) *minerConfig {
	drainConfig := drainConfig{
		Depth:       default_max_depth,
		Similarity:  default_sim,
		MaxChildren: default_max_children,
		MaxCluster:  default_max_clusters,
	}
	maskConfig := maskConfig{
		Prefix:           default_masking_prefix,
		Suffix:           default_masking_suffix,
		MaskInstructions: make([]maskInstruction, 0),
	}
	conf := minerConfig{
		Mask:  maskConfig,
		Drain: drainConfig,
	}
	for _, o := range options {
		conf = o.apply(conf)
	}
	return &conf
}

func (miner *TemplateMiner) AddLogMessage(message string) *LogMessageResponse {
	maskedMessage := miner.masker.mask(message)
	logCluster, updateType := miner.drain.addLogMessage(maskedMessage)
	return &LogMessageResponse{
		ChangeType:    updateType,
		Cluster:       logCluster,
		TemplateMined: logCluster.getTemplate(),
		ClusterCount:  len(miner.drain.idToCluster.Keys()),
	}
}

func (miner *TemplateMiner) Match(message string) *LogCluster {
	maskedMessage := miner.masker.mask(message)
	return miner.drain.match(maskedMessage, SEARCH_STRATEGY_NEVER)
}

func WithDrainDepth(depth int) minerOption {
	return minerOptionFunc(func(conf minerConfig) minerConfig {
		conf.Drain.Depth = depth
		return conf
	})
}

func WithDrainSim(sim float32) minerOption {
	return minerOptionFunc(func(conf minerConfig) minerConfig {
		conf.Drain.Similarity = sim
		return conf
	})
}

func WithDrainMaxChildren(maxChildren int) minerOption {
	return minerOptionFunc(func(conf minerConfig) minerConfig {
		conf.Drain.MaxChildren = maxChildren
		return conf
	})
}

func WithDrainMaxCluster(maxCluster int) minerOption {
	return minerOptionFunc(func(conf minerConfig) minerConfig {
		conf.Drain.MaxCluster = maxCluster
		return conf
	})
}

func WithMaskPrefix(prefix string) minerOption {
	return minerOptionFunc(func(conf minerConfig) minerConfig {
		conf.Mask.Prefix = prefix
		return conf
	})
}

func WithMaskSuffix(suffix string) minerOption {
	return minerOptionFunc(func(conf minerConfig) minerConfig {
		conf.Mask.Suffix = suffix
		return conf
	})
}

func WithMaskInsturction(pattern, maskWith string) minerOption {
	return minerOptionFunc(func(conf minerConfig) minerConfig {
		conf.Mask.MaskInstructions = append(conf.Mask.MaskInstructions, maskInstruction{
			Pattern:  pattern,
			MaskWith: maskWith,
		})
		return conf
	})
}

func (miner *TemplateMiner) Status() string {
	return miner.drain.status()
}

type minerOption interface {
	apply(minerConfig) minerConfig
}

type minerOptionFunc func(minerConfig) minerConfig

func (o minerOptionFunc) apply(conf minerConfig) minerConfig {
	return o(conf)
}
