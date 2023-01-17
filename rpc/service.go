package rpc

import (
	"go/ast"
	"log"
	"reflect"
)

type Service struct {
	name     string
	_type    reflect.Type
	instance reflect.Value
	method   map[string]*MethodType
}

func newService(instance interface{}) *Service {
	s := new(Service)
	s.instance = reflect.ValueOf(instance)
	s.name = reflect.Indirect(s.instance).Type().Name()
	// s.name = reflect.TypeOf(instance).Elem().Name()
	s._type = reflect.TypeOf(instance)
	if !ast.IsExported(s.name) {
		log.Fatalf("rcp server: %s is not a valid service name", s.name)
	}
	s.registerMethods()
	return s
}

func (s *Service) registerMethods() {
	s.method = make(map[string]*MethodType)

	for i := 0; i < s._type.NumMethod(); i++ {
		method := s._type.Method(i)
		mType := method.Type
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}
		if mType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}

		argType, replyType := mType.In(1), mType.In(2)

		s.method[method.Name] = &MethodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}

		log.Printf("rpc server: register %s.%s\n", s.name, method.Name)
	}
}

func (s *Service) call(m *MethodType, argv, replyv reflect.Value) error {

	f := m.method.Func
	returnvalue := f.Call([]reflect.Value{s.instance, argv, replyv})
	if errInter := returnvalue[0].Interface(); errInter != nil {
		return errInter.(error)
	}

	return nil
}

type MethodType struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
}

func (m *MethodType) newArgv() reflect.Value {
	var argv reflect.Value

	if m.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(m.ArgType.Elem())
	} else {
		argv = reflect.New(m.ArgType).Elem()
	}

	return argv
}

func (m *MethodType) newReplyv() reflect.Value {
	reply := reflect.New(m.ReplyType.Elem())
	switch m.ReplyType.Elem().Kind() {
	case reflect.Map:
		reply.Elem().Set(reflect.MakeMap(m.ReplyType.Elem()))
	case reflect.Slice:
		reply.Elem().Set(reflect.MakeSlice(m.ReplyType.Elem(), 0, 0))
	}

	return reply
}
