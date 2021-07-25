package structs

import (
	"github.com/knadh/koanf"
	"github.com/stretchr/testify/assert"
	"testing"
)

type A struct {
	AStuff string
}

type B struct {
	BStuff string
}

type C struct {
	A A `yaml:"a,inline"`
	B B `yaml:"b,inline"`
}

func TestInline(t *testing.T) {
	k := koanf.New(".")
	data := C{
		A: A{AStuff: "a"},
		B: B{BStuff: "b"},
	}
	assert.NoError(t, k.Load(Provider(&data, ParserTypeYaml), nil))
	assert.EqualValues(t, map[string]interface{}{"astuff": "a", "bstuff": "b"}, k.All())
}
