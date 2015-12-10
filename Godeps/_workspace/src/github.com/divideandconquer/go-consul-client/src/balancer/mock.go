package balancer

import "fmt"

type mockBalancer struct {
	services map[string]*ServiceLocation
}

func NewMockDNSBalancer(services map[string]*ServiceLocation) DNS {
	return &mockBalancer{services: services}
}

func (m *mockBalancer) FindService(serviceName string) (*ServiceLocation, error) {
	if s, ok := m.services[serviceName]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("Could not find %s", serviceName)
}
