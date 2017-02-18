package db_test

import (
	"reflect"
	"ntoolkit/component"
)

type StatelessComponentState struct {
}

type StatelessComponent struct {
	parent *component.Object
}

func (c *StatelessComponent) New() component.Component {
	return &StatelessComponent{}
}

func (c *StatelessComponent) Type() reflect.Type {
	return reflect.TypeOf(c)
}