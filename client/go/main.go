package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/krpc/krpc/client/go/krpc/pb"
	"google.golang.org/protobuf/proto"
)

func main() {

	conn, _ := net.Dial("tcp", "127.0.0.1:50000")

	connRequest := pb.ConnectionRequest{
		Type:       pb.ConnectionRequest_RPC,
		ClientName: "test",
	}

	out, err := proto.Marshal(&connRequest)
	if err != nil {
		log.Fatalln("Failed to encode event:", err)
	}
	// send to server
	fmt.Fprint(conn, string(out))
	// wait for reply
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message from server: " + message)
}
