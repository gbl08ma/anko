// +build go1.9

package packages

import (
	"reflect"
	"sync"

	"github.com/gbl08ma/anko/env"
)

func syncGo19() {
	env.PackageTypes["sync"]["Map"] = reflect.TypeOf(sync.Map{})
}
