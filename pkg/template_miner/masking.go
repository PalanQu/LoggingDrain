package templateminer

import (
	"loggingdrain/pkg/config"
	"loggingdrain/pkg/errors"
	"regexp"
)

const (
	DEFAULT_MASKING_PREFIX = "[:"
	DEFAULT_MASKING_SUFFIX = ":]"
)

type logMasker struct {
	prefix             string
	suffix             string
	nameToInstructions map[string]*logInstruction
}

type logInstruction struct {
	pattern  string
	maskWith string
	re       *regexp.Regexp
}

func newLogInstruction(maskWith, pattern string) (*logInstruction, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errors.ErrMaskPatternCompile(err)
	}
	return &logInstruction{
		pattern:  pattern,
		maskWith: maskWith,
		re:       re,
	}, nil
}

func (ins *logInstruction) mask(content, prefix, suffix string) string {
	maskStr := prefix + ins.maskWith + suffix
	return ins.re.ReplaceAllString(content, maskStr)
}

func newLogMaskerWithConfig(maskConfig config.MaskConfig) (*logMasker, error) {
	nameToInstructions := map[string]*logInstruction{}
	for _, ins := range maskConfig.MaskInstructions {
		newIns, err := newLogInstruction(ins.MaskWith, ins.Pattern)
		if err != nil {
			return nil, err
		}
		nameToInstructions[ins.MaskWith] = newIns
	}

	masker := &logMasker{
		prefix:             maskConfig.Prefix,
		suffix:             maskConfig.Suffix,
		nameToInstructions: nameToInstructions,
	}
	return masker, nil
}

func newLogMasker(prefix, suffix string) (*logMasker, error) {
	nameToInstructions := map[string]*logInstruction{}
	masker := &logMasker{
		prefix:             prefix,
		suffix:             suffix,
		nameToInstructions: nameToInstructions,
	}
	return masker, nil
}

func (mask *logMasker) addInstruction(maskWith, pattern string) error {
	ins, err := newLogInstruction(maskWith, pattern)
	if err != nil {
		return err
	}
	mask.nameToInstructions[maskWith] = ins
	return nil
}

func (mask *logMasker) mask(content string) string {
	if len(mask.nameToInstructions) == 0 {
		return content
	}
	res := content
	for _, v := range mask.nameToInstructions {
		res = v.mask(res, mask.prefix, mask.suffix)
	}
	return res
}

func (mask *logMasker) maskNames() []string {
	names := make([]string, 0, len(mask.nameToInstructions))
	for k := range mask.nameToInstructions {
		names = append(names, k)
	}
	return names
}

func (mask *logMasker) instruction(name string) *logInstruction {
	return mask.nameToInstructions[name]
}
