package packages

import (
	"os/exec"
	"reflect"

	"github.com/gbl08ma/anko/env"
)

func init() {
	env.Packages["os/exec"] = map[string]reflect.Value{
		"ErrNotFound": reflect.ValueOf(exec.ErrNotFound),
		"LookPath":    reflect.ValueOf(exec.LookPath),
		"Command":     reflect.ValueOf(exec.Command),
	}
}
