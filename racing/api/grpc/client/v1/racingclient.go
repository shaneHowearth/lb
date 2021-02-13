// Package racingclient -
package racingclient

import (
	"context"
	"log"
	"time"

	proto "github.com/shanehowearth/lb/racing/api/grpc/service/v1"
	"google.golang.org/grpc"
)

// RacingClient -
type RacingClient struct {
	Address string
}

func (rc *RacingClient) newConnection() (proto.RacingServiceClient, *grpc.ClientConn) {
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
	c, conn := rc.newConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	races, err := c.GetRaces(ctx, req)
	if err != nil {
		races = &proto.Races{}
	}
	return races, err

}