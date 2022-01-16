package pkg

import (
	"testing"
)

func TestServiceMap_HasDependent(t *testing.T) {
	sm := ServiceMap{
		"example": &Service{
			Name:     "example",
			Requires: []string{"dep1"},
		},
		"dep1": &Service{
			Name: "dep1",
		},
	}

	if sm.hasDependent("example") {
		t.Error("example should not have dependent")
	}

	if !sm.hasDependent("dep1") {
		t.Error("dep1 should have dependent")
	}
}
