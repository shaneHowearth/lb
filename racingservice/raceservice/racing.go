// Package racing -
package racing

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	grpc "github.com/shanehowearth/lb/racing/api/grpc/service/v1"
	"github.com/shanehowearth/lb/racing/repository"
)

// Server -
type Server struct {
	Store    repository.Store
	Retrieve repository.Retrieve
}

// NewRacingService -
func NewRacingService(store repository.Store, retrieve repository.Retrieve) (*Server, error) {
	if store == nil {
		return nil, fmt.Errorf("no store supplied cannot continue")
	}

	if retrieve == nil {
		return nil, fmt.Errorf("no retrieve supplied cannot continue")
	}

	return &Server{Store: store, Retrieve: retrieve}, nil
}

// https://api.ladbrokes.com.au/rest/v1/racing/?method=nextraces-category-group&count=5&include_categories=%5B%224a2788f8-e825-4d36-9894-efd4baf1cfae%22%2C%229daef0d7-bf3c-4f50-921d-8e818c60fe61%22%2C%22161d9be2-e909-4326-8c2c-35ed71fb460b%22%5D

// GetRaces -
func (s *Server) GetRaces(ctx context.Context, req *grpc.RacesRequest) (*grpc.Races, error) {

	return nil, nil
}

// CreateRaces -
func (s *Server) CreateRaces(ctx context.Context, req *grpc.Races) (*grpc.Acknowledgement, error) {
	// Handler validates the input.
	err := s.Store.CreateRaces([]repository.RaceSummary{})
	if err != nil {
		// create a unique uuid for the user to quote to tech support.
		uuid, uuiderr := uuid.NewUUID()
		if uuiderr != nil {
			// This should never happen, but if it does an alert will need to be raised immediately.
			log.Printf("Error creating uuid during article creation with context: %+v, request: %+v. error: %v", ctx, req, uuiderr)
		}
		log.Printf("Error creating article in repository: %v, code: %s", err, uuid.String())
		return &grpc.Acknowledgement{Errormessage: fmt.Sprintf("unable to create article, please quote code: %s", uuid.String())}, fmt.Errorf("unable to create article, please quote code: %s", uuid.String())
	}
	return &grpc.Acknowledgement{}, nil
}
