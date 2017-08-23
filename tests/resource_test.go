package tests

import (
	"github.com/facebookgo/inject"
	"github.com/stretchr/testify/assert"
	"github.schibsted.io/Yapo/goms/service"
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
