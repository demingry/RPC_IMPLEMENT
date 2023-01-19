package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
	"sync"
)

type Server struct {
	ServiceMap sync.Map
	ServerList *Colony
}

func (s *Server) Register(instance interface{}) {

	service := newService(instance)
	_, _ = s.ServiceMap.LoadOrStore(service.name, service)
}

func (s *Server) findServiceAndMethod(request string) (*Service, *MethodType) {

	dot := strings.LastIndex(request, ".")
	if dot < 0 {
		return nil, nil
	}

	serviceName, methodName := request[:dot], request[dot+1:]

	if srv, ok := s.ServiceMap.Load(serviceName); ok {

		service := srv.(*Service)
		if method, ok := service.method[methodName]; ok {
			return service, method
		}
	}

	return nil, nil
}

func (s *Server) StartServer(addr chan string) {

	if listener, err := net.Listen("tcp", ":8080"); err == nil {
		addr <- listener.Addr().String()
		for {
			conn, _ := listener.Accept()
			go s.serveConn(conn)
		}
	}
}

func (s *Server) serveConn(conn net.Conn) {

	var options Options

	for {
		//协商Options编解码方式
		json.NewDecoder(conn).Decode(&options)

		if options.Indicate == 0 { //初始获取codec和在服务的机器
			if len(s.ServerList.Single) == 0 {
				fmt.Fprintf(os.Stderr, "[!]There's no backend server\n")
			} else {
				options.Single = *s.dispatch()
			}
			json.NewEncoder(conn).Encode(&options)

		} else if options.Indicate == 1 { //开始接受正式数据
			json.NewEncoder(conn).Encode(&options)
			break
		} else { //不支持的Indicate码
			fmt.Fprintf(os.Stderr, "[!]Negotiate with client indicate error\n")
			return
		}
	}

	codec := InitCodecLists[options.CodecType](conn)

	for {

		var h Header
		if err := codec.ReadHeader(&h); err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				log.Println(err.Error())
			}
			return
		}

		if s, m := s.findServiceAndMethod(h.ServiceMethod); s != nil && m != nil {

			args := m.newArgv()
			reply := m.newReplyv()
			argi := args.Interface()
			if args.Type().Kind() != reflect.Ptr {
				argi = args.Addr().Interface()
			}

			if err := codec.ReadBody(argi); err != nil {
				fmt.Println(err.Error())
				return
			}

			s.call(m, args, reply)
			codec.Write(&h, reply.Interface())
		} else {
			fmt.Fprintf(os.Stderr, "cannot find service\n")
			codec.ReadBody(nil)
		}

	}

}

func (s *Server) Discover(object []*Single) {

	s.ServerList = &Colony{mu: sync.Mutex{}}
	s.ServerList.Update(object)
	// go s.ServerList.HeartBeat(10 * time.Second)
}

func (s *Server) dispatch() *Single {

	s.ServerList.mu.Lock()
	defer s.ServerList.mu.Unlock()
	if len(s.ServerList.Single) == 0 {
		return nil
	}
	single := s.ServerList.Single[s.ServerList.index]
	s.ServerList.index++
	if s.ServerList.index >= len(s.ServerList.Single) {
		s.ServerList.index = 0
	}
	return single
}
