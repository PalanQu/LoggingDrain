package templateminer

import "loggingdrain/pkg/config"

type TemplateMiner struct {
	drain  *drain
	masker *logMasker
}

type LogMessageResponse struct {
	Changetype    ClusterUpdateType
	ClusterID     int64
	TemplateMined string
	ClusterCount  int
}

func NewTemplateMinerWithConfig(config *config.Config) (*TemplateMiner, error) {
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

func NewTemplateMiner(options ...MinerOption) (*TemplateMiner, error) {
	c := newTemplateMinerConfig(options)
	return NewTemplateMinerWithConfig(c)
}

func newTemplateMinerConfig(options []MinerOption) *config.Config {
	drainConfig := config.DrainConfig{
		Depth:       DEFAULT_MAX_DEPTH,
		Similarity:  DEFAULT_SIM,
		MaxChildren: DEFAULT_MAX_CHILDREN,
		MaxCluster:  DEFAULT_MAX_CLUSTERS,
	}
	maskConfig := config.MaskConfig{
		Prefix:           DEFAULT_MASKING_PREFIX,
		Suffix:           DEFAULT_MASKING_SUFFIX,
		MaskInstructions: make([]config.MaskInstruction, 0),
	}
	conf := config.Config{
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
		Changetype:    updateType,
		ClusterID:     logCluster.id,
		TemplateMined: logCluster.getTemplate(),
		ClusterCount:  len(miner.drain.idToCluster.Keys()),
	}
}

func (miner *TemplateMiner) Match(message string) *LogCluster {
	maskedMessage := miner.masker.mask(message)
	return miner.drain.match(maskedMessage, SEARCH_STRATEGY_NEVER)
}

func WithDrainDepth(depth int) MinerOption {
	return MinerOptionFunc(func(conf config.Config) config.Config {
		conf.Drain.Depth = depth
		return conf
	})
}

func WithDrainSim(sim float32) MinerOption {
	return MinerOptionFunc(func(conf config.Config) config.Config {
		conf.Drain.Similarity = sim
		return conf
	})
}

func WithDrainMaxChildren(maxChildren int) MinerOption {
	return MinerOptionFunc(func(conf config.Config) config.Config {
		conf.Drain.MaxChildren = maxChildren
		return conf
	})
}

func WithDrainMaxCluster(maxCluster int) MinerOption {
	return MinerOptionFunc(func(conf config.Config) config.Config {
		conf.Drain.MaxCluster = maxCluster
		return conf
	})
}

func WithMaskPrefix(prefix string) MinerOption {
	return MinerOptionFunc(func(conf config.Config) config.Config {
		conf.Mask.Prefix = prefix
		return conf
	})
}

func WithMaskSuffix(suffix string) MinerOption {
	return MinerOptionFunc(func(conf config.Config) config.Config {
		conf.Mask.Suffix = suffix
		return conf
	})
}

func WithMaskInsturction(pattern, maskWith string) MinerOption {
	return MinerOptionFunc(func(conf config.Config) config.Config {
		conf.Mask.MaskInstructions = append(conf.Mask.MaskInstructions, config.MaskInstruction{
			Pattern:  pattern,
			MaskWith: maskWith,
		})
		return conf
	})
}

func (miner *TemplateMiner) Status() string {
	return miner.drain.status()
}

type MinerOption interface {
	apply(config.Config) config.Config
}

type MinerOptionFunc func(config.Config) config.Config

func (o MinerOptionFunc) apply(conf config.Config) config.Config {
	return o(conf)
}
