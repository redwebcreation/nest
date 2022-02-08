package pkg

import (
	"encoding/json"
	"os"
	"testing"
)

func TestServiceMap_BuildDependencyPlan(t *testing.T) {
	contents, err := os.ReadFile("../fixtures/dependencies.json")

	if err != nil {
		t.Fatal(err)
	}

	var datasets []struct {
		Services map[string][]string
		Sorted   [][]string
		Error    bool
	}

	err = json.Unmarshal(contents, &datasets)
	if err != nil {
		t.Fatal(err)
	}

	for _, dataset := range datasets {
		serviceMap := ServiceMap{}

		for name, dependencies := range dataset.Services {
			service := &Service{
				Name:     name,
				Requires: dependencies,
			}

			serviceMap[name] = service
		}

		sorted, err := serviceMap.BuildDependencyPlan()

		if err != nil && !dataset.Error {
			t.Errorf("Unexpected error: %s", err)
		} else if err == nil && dataset.Error {
			t.Errorf("Expected error, but got none")
		} else if err != nil && dataset.Error {
			continue
		}

		for kl, layer := range sorted {
			for ks, service := range layer {
				if service.Name != dataset.Sorted[kl][ks] {
					t.Errorf("Layer %d, service %d: expected %s, got %s", kl, ks, dataset.Sorted[kl][ks], service)
				}
			}
		}
	}
}
