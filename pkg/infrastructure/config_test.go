package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Nested struct {
	F bool `env:"LE_F"`
}

type TestConf struct {
	I int    `env:"LE_I"`
	S string `env:"LE_S"`
	F string `env:"FROM_FILE"`
	N Nested `env:"NESTED_"`
}

func TestConfigLoad(t *testing.T) {
	env := map[string]string{
		"LE_I":        "42",
		"LE_S":        "Don't panic",
		"NESTED_LE_F": "true",
	}
	// Setup environment
	for k, v := range env {
		os.Setenv(k, v)
		defer os.Unsetenv(k)
	}

	var conf TestConf
	LoadFromEnv(&conf)

	expected := TestConf{
		I: 42,
		S: "Don't panic",
		N: Nested{
			F: true,
		},
	}

	assert.Equal(t, expected, conf)
}
