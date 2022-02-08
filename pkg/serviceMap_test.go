package pkg

import (
	"sort"
	"testing"
)

func TestServiceMap_BuildDependencyPlan(t *testing.T) {
	type Set struct {
		Services map[string][]string `json:"services"`
		Expected [][]string          `json:"sorted"`
		Cyclic   bool                `json:"cyclic"`
	}

	datasets := []Set{
		{
			Services: map[string][]string{
				"example": {},
			},
			Expected: [][]string{
				{"example"},
			},
		},
		{
			Services: map[string][]string{
				"a": {"b"},
				"b": {"a"},
			},
			Cyclic: true,
		},
		{
			Services: map[string][]string{
				"a": {"b"},
				"b": {"c"},
				"c": {"a"},
			},
			Cyclic: true,
		},
		{
			Services: map[string][]string{
				"laravel": {"mysql", "redis", "elastic"},
				"redis":   {},
				"mysql":   {"fs"},
				"elastic": {"minio"},
				"minio":   {"fs"},
				"fs":      {},
			},
			Expected: [][]string{
				{"fs"},
				{"minio"},
				{"mysql", "redis", "elastic"},
				{"laravel"},
			},
		},
		{
			Services: map[string][]string{
				"example":  {"mysql"},
				"mysql":    {"fast-dfs", "logger"},
				"fast-dfs": {},
				"logger":   {},
			},
			Expected: [][]string{
				{"fast-dfs", "logger"},
				{"mysql"},
				{"example"},
			},
		},
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

		sorted, err := serviceMap.GroupServicesInLayers()

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

			sort.Strings(layer)

			for ks, service := range layer {
				if sorted[kl][ks].Name != service {
					t.Errorf("%d: expected %s, got %s (%v instead of %v)", k, service, sorted[kl][ks].Name, sorted[kl], layer)
				}
			}
		}
	}
}