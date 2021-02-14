// Package api -
package api

// Race -
type Race struct {
	Key             string
	Advertisedstart *Start
	Categoryid      string
	Meetingid       string
	Meetingname     string
	Raceform        *Form
	Raceid          string
	Racename        string
	Racenumber      int64
	Venuecountry    string
	Venueid         string
	Venuename       string
	Venuestate      string
}

// Start -
type Start struct {
	Seconds int64
}

// Form -
type Form struct {
	Additionaldata   string
	Distance         int64
	Distancetype     *DistanceType
	Distancetypeid   string
	Generated        int64
	Racecomment      string
	Racecommentalt   string
	Silkbaseurl      string
	Conditions       *TrackCondition
	Trackconditionid string
	Weather          *Weather
	Weatherid        string
}

// DistanceType -
type DistanceType struct {
	ID        string
	Name      string
	Shortname string
}

// Weather -
type Weather struct {
	Iconurl   string
	ID        string
	Name      string
	Shortname string
}

// TrackCondition -
type TrackCondition struct {
	ID        string
	Name      string
	Shortname string
}
