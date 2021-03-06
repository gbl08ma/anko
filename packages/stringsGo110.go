// +build go1.10

package packages

import (
	"reflect"
	"strings"

	"github.com/gbl08ma/anko/env"
)

func stringsGo110() {
	env.PackageTypes["strings"] = map[string]reflect.Type{
		"Builder": reflect.TypeOf(strings.Builder{}),
	}
}
