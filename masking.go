package loggingdrain

import (
	"encoding/json"
	"regexp"
)

const (
	default_masking_prefix = "[:"
	default_masking_suffix = ":]"
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

type logMaskerMarshalStruct struct {
	Prefix           string
	Suffix           string
	MaskInstructions []*logInstruction
}

type logInstructionMarshalStruct struct {
	Pattern  string
	MaskWith string
}

func (logInstruction *logInstruction) MarshalJSON() ([]byte, error) {
	marshalStruct := logInstructionMarshalStruct{
		Pattern:  logInstruction.pattern,
		MaskWith: logInstruction.maskWith,
	}
	return json.Marshal(marshalStruct)
}

func (logInstruction *logInstruction) UnmarshalJSON(data []byte) error {
	var marshalStruct logInstructionMarshalStruct
	err := json.Unmarshal(data, &marshalStruct)
	if err != nil {
		return err
	}
	re, err := regexp.Compile(marshalStruct.Pattern)
	if err != nil {
		return errMaskPatternCompile(err)
	}
	logInstruction.pattern = marshalStruct.Pattern
	logInstruction.maskWith = marshalStruct.MaskWith
	logInstruction.re = re
	return nil
}

func (logMasker *logMasker) MarshalJSON() ([]byte, error) {
	marshalStruct := logMaskerMarshalStruct{
		Prefix:           logMasker.prefix,
		Suffix:           logMasker.suffix,
		MaskInstructions: make([]*logInstruction, 0, len(logMasker.nameToInstructions)),
	}
	for _, v := range logMasker.nameToInstructions {
		marshalStruct.MaskInstructions = append(marshalStruct.MaskInstructions, v)
	}
	return json.Marshal(marshalStruct)
}

func (logMasker *logMasker) UnmarshalJSON(data []byte) error {
	var marshalStruct logMaskerMarshalStruct
	err := json.Unmarshal(data, &marshalStruct)
	if err != nil {
		return err
	}
	logMasker.prefix = marshalStruct.Prefix
	logMasker.suffix = marshalStruct.Suffix
	logMasker.nameToInstructions = map[string]*logInstruction{}
	for _, v := range marshalStruct.MaskInstructions {
		logMasker.nameToInstructions[v.maskWith] = v
	}
	return nil
}

func newLogInstruction(maskWith, pattern string) (*logInstruction, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errMaskPatternCompile(err)
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

func newLogMaskerWithConfig(maskConfig maskConfig) (*logMasker, error) {
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
