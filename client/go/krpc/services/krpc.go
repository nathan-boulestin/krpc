package services

import (
	"github.com/krpc/krpc/client/go/krpc/pb"
)

const (
	ServiceKRPC          Service   = "KRPC"
	ProcedureGetStatus   Procedure = "GetStatus"
	ProcedureGetServices Procedure = "GetServices"
)

type KRPC struct {
	connection Connection
}

func NewKRPC(conn Connection) KRPC {
	return KRPC{connection: conn}
}

func (k KRPC) GetStatus() (*pb.Status, error) {
	status := pb.Status{}
	return &status, k.connection.ProcedureCall(ServiceKRPC, ProcedureGetStatus, &status)
}

func (k KRPC) GetServices() (*pb.Services, error) {
	services := pb.Services{}
	return &services, k.connection.ProcedureCall(ServiceKRPC, ProcedureGetServices, &services)
}
