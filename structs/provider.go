// Package structs is an enhanced version of koanf's built-in structs provider.
// This version allows you to use parses such as yaml instead of fatih/structs.
// This allows you to benefit from each parser's unique struct tags when generating the map.
// For example, yaml will properly apply the "inline" tag, which is not possible with fatih/structs.
package structs

import (
	"errors"
	"gopkg.in/yaml.v3"
)

const (
	ParserTypeYaml = iota
)

type Structs struct {
	s          interface{}
	parserType int
}

// Provider returns a provider that takes a struct and a parser type
// and uses that parser to provide a confmap to koanf.
// Check the package documentation.
func Provider(s interface{}, parserType int) *Structs {
	return &Structs{s, parserType}
}

func (s *Structs) ReadBytes() ([]byte, error) {
	return nil, errors.New("structs provider does not support this method")
}

func (s *Structs) Read() (map[string]interface{}, error) {
	result := map[string]interface{}{}
	if err := s.Restructure(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Structs) Restructure(dest interface{}) error {
	switch s.parserType {
	case ParserTypeYaml:
		b, err := yaml.Marshal(s.s)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(b, dest); err != nil {
			return err
		}
		return nil
	default:
		return errors.New("unsupported parser type")
	}
}

func (s *Structs) Watch(cb func(event interface{}, err error)) error {
	return errors.New("structs provider does not support this method")
}
