// +build go1.8

package packages

import (
	"reflect"
	"time"

	"github.com/gbl08ma/anko/env"
)

func timeGo18() {
	env.Packages["time"]["Until"] = reflect.ValueOf(time.Until)
}
