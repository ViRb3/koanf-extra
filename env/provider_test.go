package env

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestInference(t *testing.T) {
	cases := []struct {
		koanfDelim string
		envDelim   string
		koanfEntry string
		envEntry   string
	}{
		{".", "_", "test.the.stuff", "TEST_THE_STUFF"},
		{".", "_", "test_the_stuff", "TEST_THE_STUFF"},
		{".", "_", "test.the_stuff", "TEST_THE_STUFF"},
		{".", "_", "test_the.stuff", "TEST_THE_STUFF"},
		//
		{"_", ".", "test.the.stuff", "TEST_THE_STUFF"},
		{"_", ".", "test_the_stuff", "TEST_THE_STUFF"},
		{"_", ".", "test.the_stuff", "TEST_THE_STUFF"},
		{"_", ".", "test_the.stuff", "TEST_THE_STUFF"},
	}
	for _, c := range cases {
		k := koanf.New(c.koanfDelim)
		os.Clearenv()
		assert.NoError(t, os.Setenv(c.envEntry, "new"), c)
		assert.NoError(t, k.Load(confmap.Provider(map[string]interface{}{c.koanfEntry: "old"}, c.koanfDelim), nil), c)
		assert.EqualValues(t, map[string]interface{}{c.koanfEntry: "old"}, k.All(), c)
		assert.NoError(t, k.Load(Provider(k, "", c.envDelim, func(s string) string {
			return strings.ToLower(s)
		}), nil), c)
		assert.EqualValues(t, map[string]interface{}{c.koanfEntry: "new"}, k.All(), c)
	}
}
