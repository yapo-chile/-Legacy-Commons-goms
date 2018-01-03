package tests

import (
	"github.schibsted.io/Yapo/goms/service"
	"gopkg.in/facebookgo/inject.v0"
	"gopkg.in/stretchr/testify.v1/assert"
	"testing"
)

/* Example test case */
func TestResource(t *testing.T) {
	name := "yo"
	usage := "no fui"

	service.SetupInject(
		&inject.Object{Value: &service.Resource{Name: name, Usage: usage}},
	)
	resource := service.Inject("Resource").(*service.Resource)

	assert.Equal(t, resource.SumLength(), len(name)+len(usage))
}

func TestConf(t *testing.T) {
	svc := service.RuntimeConfig{
		Host:    "http://juana.la/iguana",
		Port:    123123,
		PidFile: "Woopsie",
	}

	service.SetConf(&service.Config{
		Runtime: svc,
	})

	conf := service.GetConfig()
	assert.Equal(t, conf.Runtime, svc)
}
