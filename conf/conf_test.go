package conf

import (
	"gopkg.in/stretchr/testify.v1/assert"
	"testing"
)

func TestLoadConfNoFile(t *testing.T) {
	conf, err := LoadConf("testdata/nosuchfile.json")
	assert.Nil(t, conf)
	assert.Error(t, err)
}

func TestLoadConfBadFormat(t *testing.T) {
	conf, err := LoadConf("testdata/badformat.json")
	assert.Nil(t, conf)
	assert.Error(t, err)
}

func TestLoadConfOk(t *testing.T) {
	conf, err := LoadConf("testdata/ok.json")
	expected := &Config{
		Runtime: RuntimeConfig{
			Host:    "somehost",
			Port:    12312,
			PidFile: "lepid",
		},
	}
	assert.Equal(t, expected, conf)
	assert.Nil(t, err)
}

func TestConf(t *testing.T) {
	svc := RuntimeConfig{
		Host:    "http://juana.la/iguana",
		Port:    123123,
		PidFile: "Woopsie",
	}

	Set(&Config{
		Runtime: svc,
	})

	conf := Get()
	assert.Equal(t, conf.Runtime, svc)
}
