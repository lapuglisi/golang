package coredb

import (
	"plugin"
	"reflect"
)

// ADTypesByFile .....
type ADTypesByFile map[string]reflect.Type

var coredbTypes ADTypesByFile

func init() {
	coredbTypes = make(map[string]reflect.Type)

	p, err := plugin.Open("dbsoler_types.so")
	if err != nil {
		panic("Miou porque nao achou plugin.")
	}

	m, err := p.Lookup("InitTypes")
	if err != nil {
		panic("Miou porque nao achou InitTypes.")
	}

	// Calls InitTypes from within the plugin
	m.(func(map[string]reflect.Type))(coredbTypes)
}

// GetStructType ......
func GetStructType(name string) reflect.Type {
	value, ok := coredbTypes[name]

	if !ok {
		return nil
	}

	return value
}
