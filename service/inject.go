package service

import (
	"github.com/facebookgo/inject"
	"reflect"
)

/*
 * Dependency Injection implementation:
 *
 * Fill the injectData struct below with pointer declarations to
 * every resource your Handlers may need to inject. It works better
 * when you use interfaces than classes.
 *
 * Then, ensure that you provide a concrete object instance for each
 * of the listed interfaces. Those will be used to fill up your objects
 * whenever such interface is required by another object. Use the
 * SetupInject function for that.
 *
 * Finally, you can use `Inject("Resource")` on your Handler code to
 * receive a ready to use pointer to one of the inflated structs.
 *
 * Note: The caller is responsible for casting the received resource
 * to the correct interface. This will not only allow them to use the
 * object properly, but will result on predictable errors whenerver
 * expectations are not met.
 */

var injectdata injectData

type injectData struct {
	/* Add your resources here */
	Resource *Resource `inject:""`
}

/*
 * SetupInject must be invoked very early on your service lifetime
 * and be provided with an instance of EVERY injectable resource
 * in the code base. Running "git grep inject:" should help you
 * find a list of everything that needs to be defined.
 */
func SetupInject(values ...*inject.Object) error {
	var graph inject.Graph
	values = append(values, &inject.Object{Value: &injectdata})
	err := graph.Provide(values...)
	if err != nil {
		return err
	}
	return graph.Populate()
}

/*
 * Inject returns a resource pointer to a ready to use instance of
 * the requested resource name. Names should be provided as they
 * appear on the injectData struct or very ugly errors will occur.
 * The caller should cast the received resource to the interface
 * they expected to get. This will result in predictable conversion
 * errors when the received pointer is not what the caller expected
 */
func Inject(resource string) interface{} {
	instance := reflect.ValueOf(injectdata).FieldByName(resource)
	return instance.Interface()
}
