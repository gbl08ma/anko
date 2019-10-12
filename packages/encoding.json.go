package packages

import (
	"reflect"
	"encoding/json"

	"github.com/gbl08ma/anko/env"
)

func init() {
	env.Packages["encoding/json"] = map[string]reflect.Value{
		"Marshal":   reflect.ValueOf(json.Marshal),
		"Unmarshal": reflect.ValueOf(json.Unmarshal),
	}
}
