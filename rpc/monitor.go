package rpc

import (
	"fmt"
	"net/http"
)

type Monitor struct {
	server *Server
}

func (m Monitor) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	//服务调用数量
	m.server.ServiceMap.Range(func(key, value any) bool {
		service := value.(*Service)
		fmt.Fprintf(w, "Service: %s\n", service.name)
		for name, function := range service.method {
			fmt.Fprintf(w, "Method: %s invoked %d\n", name, function.invoked)
		}

		return true
	})

	//服务主机数量
	m.server.ServerList.mu.Lock()
	fmt.Fprintf(w, "Current servers: %d\n", len(m.server.ServerList.Single))
	m.server.ServerList.mu.Unlock()
}

func (m *Monitor) StartMonitor() {

	http.Handle("/monitor", m)
	http.ListenAndServe(":8081", nil)
}
