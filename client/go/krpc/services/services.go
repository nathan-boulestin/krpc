package services

import "google.golang.org/protobuf/reflect/protoreflect"

type Service string

type Procedure string

type Connection interface {
	ProcedureCall(s Service, p Procedure, response protoreflect.ProtoMessage) error
}
