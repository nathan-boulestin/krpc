package main

import (
	"fmt"
	"log"

	"github.com/krpc/krpc/client/go/krpc"
	"github.com/krpc/krpc/client/go/krpc/services"
)

func main() {

	conn, err := krpc.NewConnection("test", krpc.DefaultHost, krpc.DefaultPort, krpc.NoStream)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Printf("Connected with client %s\n", conn.Name)

	krpcService := services.NewKRPC(conn)
	status, err := krpcService.GetStatus()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("KRPC version: %s\n", status.Version)

	services, err := krpcService.GetServices()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("services: %v\n", services.String())
}
