// Code generated by ""fitask" -type=Network"; DO NOT EDIT

package gcetasks

import (
	"encoding/json"

	"k8s.io/kops/upup/pkg/fi"
)

// Network

// JSON marshalling boilerplate
type realNetwork Network

func (o *Network) UnmarshalJSON(data []byte) error {
	var jsonName string
	if err := json.Unmarshal(data, &jsonName); err == nil {
		o.Name = &jsonName
		return nil
	}

	var r realNetwork
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	*o = Network(r)
	return nil
}

var _ fi.HasName = &Network{}

func (e *Network) GetName() *string {
	return e.Name
}

func (e *Network) SetName(name string) {
	e.Name = &name
}

func (e *Network) String() string {
	return fi.TaskAsString(e)
}
