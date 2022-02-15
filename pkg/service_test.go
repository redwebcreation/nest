package pkg

import (
	"sort"
	"testing"
)

func TestService_ApplyDefaults(t *testing.T) {
	service := &Service{
		Hosts: []string{
			"~example.com",
			"git.example.com",
		},
	}

	service.ApplyDefaults("example")

	if len(service.Hosts) != 3 {
		t.Errorf("Expected 3 hosts, got %d", len(service.Hosts))

		if service.Hosts[0] != "example.com" && service.Hosts[1] != "www.example.com" && service.Hosts[2] != "git.example.com" {
			t.Errorf("Expected example.com, www.example.com, git.example.com, got %s, %s, %s", service.Hosts[0], service.Hosts[1], service.Hosts[2])
		}
	}

	if service.ListeningOn != "80" {
		t.Errorf("Expected s.ListeningOn to default to 80, got %s", service.ListeningOn)
	}

	service = &Service{
		ListeningOn: ":port",
	}

	service.ApplyDefaults("example")

	if service.ListeningOn != "port" {
		t.Errorf("Expected s.ListeningOn to trim leading colon, got %s", service.ListeningOn)
	}

	service = &Service{
		ListeningOn: "443",
	}

	service.ApplyDefaults("example")

	if service.ListeningOn != "443" {
		t.Errorf("Expected s.ListeningOn to be unchanged, got %s", service.ListeningOn)
	}
}

func TestServiceMap_BuildDependencyPlan(t *testing.T) {
	type Test struct {
		Services map[string][]string `json:"services"`
		Expected [][]string          `json:"sorted"`
		Cyclic   bool                `json:"cyclic"`
	}

	tests := []Test{
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

	for k, test := range tests {
		serviceMap := ServiceMap{}

		for name, dependencies := range test.Services {
			service := &Service{
				Name:     name,
				Requires: dependencies,
			}

			serviceMap[name] = service
		}

		sorted, err := serviceMap.GroupInLayers()

		if test.Cyclic {
			if err == nil {
				t.Errorf("%d: expected error, got nil", k)
			}

			continue
		}

		if err != nil {
			t.Errorf("%d: unexpected error: %s", k, err)
		}

		if len(sorted) != len(test.Expected) {
			t.Fatalf("%d: expected %d services, got %d", k, len(test.Expected), len(sorted))
		}

		for kl, layer := range test.Expected {
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
