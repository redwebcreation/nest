package pkg

import "testing"

func TestService_Normalize(t *testing.T) {
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
