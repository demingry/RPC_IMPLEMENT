package rpc

import (
	"fmt"
	"net/http"
)

type Monitor struct {
	server *Server
}

func (m Monitor) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.server.ServiceMap.Range(func(key, value any) bool {
		service := value.(*Service)
		fmt.Fprintf(w, "Service: %s\n", service.name)
		for name, function := range service.method {
			fmt.Fprintf(w, "Method: %s invoked %d\n", name, function.invoked)
		}

		return true
	})
}

func (m *Monitor) StartMonitor() {

	http.Handle("/monitor", m)
	http.ListenAndServe(":8081", nil)
}
