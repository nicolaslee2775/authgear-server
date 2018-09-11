package inject

import (
	"context"
	"reflect"

	"github.com/skygeario/skygear-server/pkg/core/config"
	"github.com/skygeario/skygear-server/pkg/core/handler"
)

func DefaultInject(
	i interface{},
	dependencyGraph handler.ProviderGraph,
	ctx context.Context,
	configuration config.TenantConfiguration,
) {
	injectDependency(i, dependencyGraph, ctx, configuration)
}

func injectDependency(
	i interface{},
	dependencyGraph handler.ProviderGraph,
	ctx context.Context,
	configuration config.TenantConfiguration,
) {
	t := reflect.TypeOf(i).Elem()
	v := reflect.ValueOf(i).Elem()

	numField := t.NumField()
	for i := 0; i < numField; i++ {
		dependencyName := t.Field(i).Tag.Get("dependency")
		if dependencyName == "" {
			continue
		}

		field := v.Field(i)
		dependency := dependencyGraph.Provide(dependencyName, ctx, configuration)
		field.Set(reflect.ValueOf(dependency))
	}
}