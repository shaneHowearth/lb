package main

import (
	"context"
	"log"
	"net"
	"os"

	// racing "github.com/shanehowearth/lb/racing/raceservice"
	// grpcservice "github.com/shanehowearth/lb/racingservice/api/grpc/service/v1"
	grpcservice "github.com/shanehowearth/lb/racingservice/api/grpc/service/v1"
	racing "github.com/shanehowearth/lb/racingservice/raceservice"

	// repo "github.com/shanehowearth/lb/racingservice/repository/postgres/v1"
	"github.com/shanehowearth/lb/racingservice/repository/postgres/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

var ss *racing.Server

func main() {

	ok := false

	// Datastore
	ds := new(postgres.Postgres)
	ds.Retry = 1

	ds.URI, ok = os.LookupEnv("DBURI")
	if !ok {
		log.Fatalf("DBURI is not set")
	}

	// Racing Service
	var err error
	ss, err = racing.NewRacingService(ds, ds)
	if err != nil {
		log.Fatalf("Failed to obtain New Racing Service with error %v", err)
	}

	// gRPC service
	portNum := os.Getenv("PORT_NUM")
	lis, err := net.Listen("tcp", "0.0.0.0:"+portNum)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	grpcservice.RegisterRacingServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateRaces -
func (s *server) CreateRaces(ctx context.Context, races *grpcservice.Races) (*grpcservice.Acknowledgement, error) {

	st, err := ss.CreateRaces(ctx, races)
	if err != nil {
		return nil, err
	}
	return st, nil
}

// GetRaces -
func (s *server) GetRaces(ctx context.Context, req *grpcservice.RacesRequest) (*grpcservice.Races, error) {
	log.Printf("Server GetRaces called with %v", req)
	races, err := ss.GetRaces(ctx, req)
	if err != nil {
		log.Printf("Server returning error %v", err)
		return nil, err
	}
	log.Printf("Server returning races %v", races)
	return races, nil
}
