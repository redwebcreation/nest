package common

import "testing"

func TestService_Accepts(t *testing.T) {
	service := Service{
		Hosts: []string{"example.com", "*.example.com", "app.*.example.com", "*.*.example.com"},
	}

	if !service.Accepts("example.com") {
		t.Error("Service should accept example.com")
	}

	if !service.Accepts("www.example.com") {
		t.Error("Service should accept www.example.com")
	}

	if !service.Accepts("app.www.example.com") {
		t.Error("Service should accept app.customer1.example.com")
	}

	if service.Accepts("") {
		t.Error("Service should not accept empty string")
	}

}
