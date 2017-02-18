package db_test

import (
	"reflect"
	"ntoolkit/component"
)

type StatefulComponentState struct {
	Name  string
	Value string
}

type StatefulComponent struct {
	parent *component.Object
	state  StatefulComponentState
}

func (c *StatefulComponent) New() component.Component {
	return &StatefulComponent{}
}

func (c *StatefulComponent) Type() reflect.Type {
	return reflect.TypeOf(c)
}

func (c *StatefulComponent) Serialize() (interface{}, error) {
	return component.SerializeState(c.state)
}

func (c *StatefulComponent) Deserialize(raw interface{}) error {
	return component.DeserializeState(&c.state, raw)
}
