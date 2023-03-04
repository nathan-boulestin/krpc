package krpc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/krpc/krpc/client/go/krpc/pb"
	"github.com/krpc/krpc/client/go/krpc/services"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	DefaultHost       string = "127.0.0.1"
	DefaultPort       string = "50000"
	DefaultStreamPort string = "50001"
	NoStream          string = ""
)

type Connection struct {
	Name     string
	clientId []byte
	conn     net.Conn
	stream   net.Conn
}

func (c Connection) Close() {
	c.conn.Close()
}

func NewConnection(name, ip, port, streamPort string) (Connection, error) {
	c := Connection{
		Name: name,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return Connection{}, err
	}
	c.conn = conn

	connectionRequest := pb.ConnectionRequest{
		Type:       pb.ConnectionRequest_RPC,
		ClientName: name,
	}
	connectionResponse := pb.ConnectionResponse{}

	err = c.communicate(&connectionRequest, &connectionResponse)
	if err != nil {
		return Connection{}, err
	}
	if connectionResponse.Status != pb.ConnectionResponse_OK {
		return Connection{}, fmt.Errorf("failled connecting: status:%s, message:%s", connectionResponse.Status.String(), connectionResponse.Message)
	}
	c.clientId = connectionResponse.ClientIdentifier

	if streamPort != NoStream {
		stream, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, streamPort))
		if err != nil {
			return Connection{}, err
		}
		c.stream = stream
		streamConnectionRequest := pb.ConnectionRequest{
			Type:             pb.ConnectionRequest_STREAM,
			ClientIdentifier: c.clientId,
		}
		streamConnectionResponse := pb.ConnectionResponse{}
		err = c.communicateStream(&streamConnectionRequest, &streamConnectionResponse)
		if err != nil {
			return Connection{}, err
		}
		if streamConnectionResponse.Status != pb.ConnectionResponse_OK {
			return Connection{}, fmt.Errorf("failled connecting: status:%s, message:%s", streamConnectionResponse.Status.String(), streamConnectionResponse.Message)
		}
	}

	return c, nil
}

func (c Connection) communicate(message, response protoreflect.ProtoMessage) error {
	messageBytes, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	messageSizeBytes := protowire.AppendVarint(nil, uint64(proto.Size(message)))
	_, err = c.conn.Write(append(messageSizeBytes, messageBytes...))
	if err != nil {
		return err
	}

	var responseSize uint64
	var n int
	responseSizeBytes := []byte{}
	for {
		nextByte := make([]byte, 1)
		_, err = c.conn.Read(nextByte)
		if err != nil {
			return err
		}
		responseSizeBytes = append(responseSizeBytes, nextByte...)
		responseSize, n = protowire.ConsumeVarint(responseSizeBytes)
		if n < 0 {
			err := protowire.ParseError(n)
			if errors.Is(err, io.ErrUnexpectedEOF) {
				continue
			}
			return err
		}
		break
	}
	responseBytes := make([]byte, responseSize)
	_, err = c.conn.Read(responseBytes)
	if err != nil {
		return err
	}
	return proto.Unmarshal(responseBytes, response)
}

func (c Connection) communicateStream(message, response protoreflect.ProtoMessage) error {
	messageBytes, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	messageSizeBytes := protowire.AppendVarint(nil, uint64(proto.Size(message)))
	_, err = c.stream.Write(append(messageSizeBytes, messageBytes...))
	if err != nil {
		return err
	}

	reader := bufio.NewReader(c.stream)
	responseSizeBytes := make([]byte, 1)
	_, err = reader.Read(responseSizeBytes)
	if err != nil {
		return err
	}
	responseSize, _ := protowire.ConsumeVarint(responseSizeBytes)

	responseBytes := make([]byte, responseSize)
	_, err = reader.Read(responseBytes)
	if err != nil {
		return err
	}
	return proto.Unmarshal(responseBytes, response)
}

func (c Connection) ProcedureCall(s services.Service, p services.Procedure, response protoreflect.ProtoMessage) error {
	call := pb.ProcedureCall{
		Service:   string(s),
		Procedure: string(p),
	}
	request := pb.Request{
		Calls: []*pb.ProcedureCall{&call},
	}
	resp := pb.Response{}
	if err := c.communicate(&request, &resp); err != nil {
		return err
	}
	if resp.Error != nil {
		return fmt.Errorf("request failed: %s", resp.Error.String())
	}
	fmt.Printf("len(resp.Results): %v\n", len(resp.Results))
	if len(resp.Results) == 1 {
		if resp.Results[0].Error != nil {
			return fmt.Errorf("request failed: %s", resp.Results[0].Error.String())
		}
	}
	if response != nil {
		if err := proto.Unmarshal(resp.Results[0].Value, response); err != nil {
			return err
		}
	}
	return nil
}
