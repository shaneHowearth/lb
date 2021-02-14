// Package racingclient -
package racingclient

import (
	"context"
	"log"
	"time"

	proto "github.com/shanehowearth/lb/racingservice/api/grpc/service/v1"
	"google.golang.org/grpc"
)

// RacingClient -
type RacingClient struct {
	Address string
}

func (rc *RacingClient) newConnection() (proto.RacingServiceClient, *grpc.ClientConn) {
	log.Printf("New connection address %s", rc.Address)
	// Set up a connection to the server.
	conn, err := grpc.Dial(rc.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return proto.NewRacingServiceClient(conn), conn
}

// CreateRaces -
func (rc *RacingClient) CreateRaces(races *proto.Races) *proto.Acknowledgement {
	c, conn := rc.newConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ack, err := c.CreateRaces(ctx, races)
	if err != nil {
		ack = &proto.Acknowledgement{Errormessage: err.Error()}
	}
	return ack
}

// GetRaces -
func (rc *RacingClient) GetRaces(req *proto.RacesRequest) (*proto.Races, error) {
	log.Printf("RacingClient GetRaces received req: %v", req)
	c, conn := rc.newConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// log.Printf("RacingClient GetRaces received req: %v", req)

	races, err := c.GetRaces(ctx, req)
	if err != nil {
		log.Printf("RacingClient GetRaces returning error: %v", err)
		races = &proto.Races{}
	}
	log.Printf("RacingClient GetRaces returning error: %v", err)
	return races, err

}
