package ws

// Service interface to define behavior for services
type Service interface {
	Process(input map[string]interface{}) (map[string]interface{}, error)
}

// Service registry to hold service constructors
var serviceRegistry = make(map[string]func() Service)

// RegisterService allows services to register themselves
func RegisterService(name string, constructor func() Service) {
	serviceRegistry[name] = constructor
}

// GetService retrieves a new instance of a service by name
func GetService(name string) (Service, bool) {
	constructor, exists := serviceRegistry[name]
	if !exists {
		return nil, false
	}
	return constructor(), true
}

// ListServices returns a slice of the names of all registered services
func ListServices() []string {
	names := make([]string, 0, len(serviceRegistry))
	for name := range serviceRegistry {
		names = append(names, name)
	}
	return names
}
