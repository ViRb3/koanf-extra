// Package env is an enhanced version of koanf's built-in env provider.
// This version is able to parse environment variables that are ambiguous due to their delimiter.
// By convention, environment variables use an underscore '_' as delimiter. YAML entries also use the same delimiter.
// This creates an ambiguity problem. For an example, take the following structure:
//   test:
//     the_stuff:
// The corresponding environment variable would ideally be:
//    TEST_THE_STUFF
// However, a parser cannot know whether THE and STUFF are two nested entries or one entry with the name 'THE_STUFF'.
// This provider will infer the correct name for each environment variable based on the exising data stored in koanf.Koanf.
package env

import (
	"errors"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/maps"
	"github.com/knadh/koanf/providers/env"
	"reflect"
	"strings"
)

type Env struct {
	k        *koanf.Koanf
	envDelim string
	orig     *env.Env
}

// Provider returns a provider that will infer the correct name for each environment variable based on the exising data stored in koanf.Koanf.
// Check the package documentation.
func Provider(k *koanf.Koanf, prefix, delim string, cb func(s string) string) *Env {
	return &Env{k, delim, env.Provider(prefix, delim, cb)}
}

func (s *Env) ReadBytes() ([]byte, error) {
	return nil, errors.New("env provider does not support this method")
}

func (s *Env) Read() (map[string]interface{}, error) {
	result, err := s.orig.Read()
	if err != nil {
		return nil, err
	}
	result, _ = maps.Flatten(result, []string{}, s.k.Delim())
	result = s.inferDelimiters(result)
	result = maps.Unflatten(result, s.k.Delim())
	return result, nil
}

func (s *Env) Watch(cb func(event interface{}, err error)) error {
	return errors.New("env provider does not support this method")
}

// inferDelimiters will take a configuration map with both properly separated entries (test.the_stuff) and
// greedily separated entries (test.the.stuff), and given the map separation delimiter, will transfer the
// values of all greedily separated entries onto the properly separated entries. This process effectively "fixes" the
// ambiguity problem and allows the configuration map to be used as usual.
// Check the package documentation.
func (s *Env) inferDelimiters(newData map[string]interface{}) map[string]interface{} {
	existingData := s.k.All()
	result := map[string]interface{}{}
	for name := range existingData {
		for name2, val2 := range newData {
			splitFunc := func(r rune) bool {
				return string(r) == s.envDelim || string(r) == s.k.Delim()
			}
			split := strings.FieldsFunc(name, splitFunc)
			split2 := strings.FieldsFunc(name2, splitFunc)
			if reflect.DeepEqual(split, split2) {
				result[name] = val2
			}
		}
	}
	return result
}
