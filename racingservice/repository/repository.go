// Package repository -
package repository

// Provide interfaces for interaction with the repository. Having an interface
// allows the implementation(s) to be changed at will, without the business
// logic requiring any change (the BL won't know what repository it is
// interacting with, the Dependency Injection will).

// Store - Store a race summary
type Store interface {
	CreateRaces([]RaceSummary) error
}

// Retrieve - Retrieve a previously stored race summary
type Retrieve interface {
	GetRaces(string, int64, []string) ([]RaceSummary, error)
}

// RaceSummary - DTO
type RaceSummary struct {
	ID              string
	AdvertisedStart struct {
		Seconds int
	}
	CategoryID  string
	MeetingID   string
	MeetingName string
	RaceForm    struct {
		AdditionalData string
		Distance       int
		DistanceType   struct {
			ID        string
			Name      string
			ShortName string
		}
		DistanceTypeID         string
		Generated              int
		RaceComment            string
		RaceCommentAlternative string
		SilkBaseURL            string
		TrackCondition         struct {
			ID        string
			Name      string
			ShortName string
		}
		TrackConditionID string
		Weather          struct {
			IconURI   string
			ID        string
			Name      string
			ShortName string
		}
		WeatherID string
	}
	RaceFormID   string
	RaceID       string
	RaceName     string
	RaceNumber   int
	VenueCountry string
	VenueID      string
	VenueName    string
	VenueState   string
}
