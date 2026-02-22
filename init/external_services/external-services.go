package external_services

import (
	"net/http"
	"sync"
	"time"

	"github.com/gustavo000/goLibGustavo/models/properties"
	"github.com/gustavo000/goLibGustavo/models/rest"
)

func GetEndpointsOfServices(environment string) {
	group := &sync.WaitGroup{}
	var mutex sync.RWMutex
	group.Add(len(properties.GetProperty().External.Services))
	for i, service := range properties.GetProperty().External.Services {
		go func(i int, service *rest.Service) {
			defer group.Done()
			if service.Timeout == 0 {
				service.Timeout = 10
			}
			if service.Client == nil {
				timeout := time.Second * time.Duration(service.Timeout)
				service.Client = &http.Client{Timeout: timeout}
			}
			/*			if !strings.Contains(strings.ToLower(service.Layer), "external") {
						service = endpoints.GetEndpoints(service, environment)
					}*/
			mutex.Lock()
			properties.GetProperty().External.Services[i] = service
			mutex.Unlock()
		}(i, service)
	}
	group.Wait()
	return
}
