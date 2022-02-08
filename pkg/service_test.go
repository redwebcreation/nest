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
		Services map[string][]string `json:"services"`
		Expected [][]string          `json:"sorted"`
		Cyclic   bool                `json:"cyclic"`
	}

	err = json.Unmarshal(contents, &datasets)
	if err != nil {
		t.Fatal(err)
	}

	for k, dataset := range datasets {
		serviceMap := ServiceMap{}

		for name, dependencies := range dataset.Services {
			service := &Service{
				Name:     name,
				Requires: dependencies,
			}

			serviceMap[name] = service
		}

		sorted, err := serviceMap.BuildDependencyPlan()

		if dataset.Cyclic {
			if err == nil {
				t.Errorf("%d: expected error, got nil", k)
			}

			continue
		}

		if err != nil {
			t.Errorf("%d: unexpected error: %s", k, err)
		}

		if len(sorted) != len(dataset.Expected) {
			t.Fatalf("%d: expected %d services, got %d", k, len(dataset.Expected), len(sorted))
		}

		for kl, layer := range dataset.Expected {
			if len(layer) != len(sorted[kl]) {
				t.Fatalf("%d: expected %d services in layer %d, got %d", k, len(layer), kl, len(sorted[kl]))
			}

			for ks, service := range layer {
				if sorted[kl][ks].Name != service {
					t.Errorf("%d: expected %s, got %s", k, service, sorted[kl][ks].Name)
				}
			}
		}
	}
}
