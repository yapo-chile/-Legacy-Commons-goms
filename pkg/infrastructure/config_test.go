package infrastructure

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Nested struct {
	F bool `env:"LE_F"`
}

type TestConf struct {
	I  int    `env:"LE_I"`
	S  string `env:"LE_S"`
	F  string `env:"FROM"`
	N  Nested `env:"NESTED_"`
	D  string `env:"DEF" envDefault:"default_conf"`
	OF string `env:"OTHERFILE"`
}

func createFile(prefix string) string {
	content := []byte("temp_data")
	tmpfile, err := ioutil.TempFile("", prefix)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	return tmpfile.Name() // remember clean up
}

func TestConfigLoad(t *testing.T) {
	dummyFile := createFile("lala")
	defer os.Remove(dummyFile)
	env := map[string]string{
		"LE_I":           "42",
		"LE_S":           "Don't panic",
		"NESTED_LE_F":    "true",
		"FROM_FILE":      dummyFile,
		"OTHERFILE_FILE": "/fileMissing.txt",
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
		F: "temp_data",
		N: Nested{
			F: true,
		},
		D: "default_conf",
	}

	assert.Equal(t, expected, conf)
}
