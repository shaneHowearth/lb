// Package postgres -
package postgres

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shanehowearth/lb/racingservice/repository"
)

// Postgres -
type Postgres struct {
	pool  *pgxpool.Pool
	Retry int
	URI   string
}

// Expose sql.Open to testing within this package
// var sqlOpen = sql.Open

// Connect - Create the connection to the database
func (p *Postgres) Connect() (err error) {
	// Retry MUST be >= 1
	if p.Retry == 0 {
		log.Print("Cannot use a Retry of zero, this process will to default retry to 5")
		p.Retry = 5
	}
	if p.URI == "" {
		log.Print("no Postgres URI configured")
		return fmt.Errorf("no Postgres URI configured")
	}

	// Infinite loop
	// Keep trying forever
	for {
		for i := 0; i < p.Retry; i++ {
			p.pool, err = pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))

			if err != nil {
				log.Printf("Unable to connect to database: %v\n", err)
			}
			if p.pool != nil {
				return nil
			}
			time.Sleep(1 * time.Second)
		}

		backoff := time.Duration(p.Retry*rand.Intn(10)) * time.Second
		log.Printf("ALERT: Trouble connecting to Postgres, error: %v, going to re-enter retry loop in %s seconds", err, backoff.String())
		time.Sleep(backoff)
	}
}

// CreateRaces -
func (p *Postgres) CreateRaces(races []repository.RaceSummary) error {
	if p.pool == nil {
		perr := p.Connect()
		if perr != nil {
			// should never get here
			log.Fatalf("unable to connect to postgres server with error: %v", perr)
		}
	}

	// Begin Transaction
	tx, err := p.pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("unable to create transaction with error %w", err)
	}
	// Rollback is fine even on success (There'll be nothing to rollback after
	// the commit)
	defer func() {
		txErr := tx.Rollback(context.Background())
		if txErr != nil {
			log.Printf("Error generated when rolling back transaction %v", txErr)
		}
	}()

	for _, race := range races {
		// Weather Conditions
		err := tx.QueryRow(context.Background(), `INSERT INTO weather_conditions(icon_uri, id, name, short_name) VALUES($1, $2, $3, $4)`, race.RaceForm.Weather.IconURI, race.RaceForm.Weather.ID, race.RaceForm.Weather.Name, race.RaceForm.Weather.ShortName)
		if err != nil {
			return fmt.Errorf("unable to insert into weather_conditions with error %w", err)
		}

		// Track conditions
		err = tx.QueryRow(context.Background(), `INSERT INTO track_conditions(id, name, short_name) VALUES($1, $2, $3)`, race.RaceForm.TrackCondition.ID, race.RaceForm.TrackCondition.Name, race.RaceForm.TrackCondition.ShortName)
		if err != nil {
			return fmt.Errorf("unable to insert into track_condition with erro %w", err)
		}

		// Distance types
		err = tx.QueryRow(context.Background(), `INSERT INTO distance_types(id, name, short_name) VALUES($1, $2, $3)`, race.RaceForm.DistanceType.ID, race.RaceForm.DistanceType.Name, race.RaceForm.DistanceType.ShortName)
		if err != nil {
			return fmt.Errorf("unable to insert into distance_types with erro %w", err)
		}

		// Race Form
		err = tx.QueryRow(context.Background(), `INSERT INTO race_forms(id, addtional_data, distance, distance_type_id, generated, race_comment, race_comment_alternative, silk_base_url, track_condition_id, weather_id) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, race.RaceFormID, race.RaceForm.AdditionalData, race.RaceForm.Distance, race.RaceForm.DistanceTypeID, race.RaceForm.Generated, race.RaceForm.RaceComment, race.RaceForm.RaceCommentAlternative, race.RaceForm.SilkBaseURL, race.RaceForm.TrackConditionID, race.RaceForm.WeatherID)
		if err != nil {
			return fmt.Errorf("unable to insert into race_form with erro %w", err)
		}

		// Race Summaries
		err = tx.QueryRow(context.Background(), `INSERT INTO race_summaries(id, advertised_start, category_id, meeting_id, meeting_name, race_form, race_id, race_name, race_number, venue_country, venue_id, venue_name, venue_state), VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`, race.ID, race.AdvertisedStart.Seconds, race.CategoryID, race.MeetingID, race.MeetingName, race.RaceFormID, race.RaceID, race.RaceName, race.RaceNumber, race.VenueCountry, race.VenueID, race.VenueName, race.VenueState)

		if err != nil {
			return fmt.Errorf("unable to insert into race_summaries with erro %w", err)
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("unable to commit transaction with error %w", err)
	}
	return nil
}

// GetNextRaces -
func (p *Postgres) GetNextRaces(count int64, races []string) ([]repository.RaceSummary, error) {
	log.Printf("Postgres GetRaces called with %v", races)
	if p.pool == nil {
		perr := p.Connect()
		if perr != nil {
			// should never get here
			log.Fatalf("unable to connect to postgres server with error: %v", perr)
		}
	}
	results := []repository.RaceSummary{}

	var summary []struct {
		ID           string `db:"id"`
		Seconds      int    `db:"advertised_start"`
		CategoryID   string `db:"category_id"`
		MeetingID    string `db:"meeting_id"`
		MeetingName  string `db:"meeting_name"`
		RaceFormID   string `db:"race_form"`
		RaceID       string `db:"race_id"`
		RaceName     string `db:"race_name"`
		RaceNumber   int    `db:"race_number"`
		VenueCountry string `db:"venue_country"`
		VenueID      string `db:"venue_id"`
		VenueName    string `db:"venue_name"`
		VenueState   string `db:"venue_state"`
	}
	// Race summary
	err := p.pool.QueryRow(context.Background(), "SELECT id, advertised_start, category_id, meeting_id, meeting_name, race_form, race_id, race_name, race_number, venue_country, venue_id, venue_name, venue_state FROM race_summaries WHERE id = ANY ($1) LIMIT $2", races, count).Scan(&summary)
	if err != nil {
		return []repository.RaceSummary{}, fmt.Errorf("unable to fetch summary with id %s failed: %w", races, err)
	}

	for idx := range summary {
		// Raceform
		var raceform struct {
			ID                     string `db:"id"`
			AdditionalData         string `db:"additional_data"`
			Distance               int    `db:"distance"`
			DistanceTypeID         string `db:"distance_type_id"`
			Generated              int    `db:"generated"`
			RaceComment            string `db:"race_comment"`
			RaceCommentAlternative string `db:"race_comment_alternative"`
			SilkBaseURL            string `db:"silk_base_url"`
			TrackConditionID       string `db:"track_condition_id"`
			WeatherID              string `db:"weather_id"`
		}
		err = p.pool.QueryRow(context.Background(), "SELECT id, additional_data, distance, distance_type_id, generated, race_comment, race_comment_alternative, silk_base_url, track_condition_id, weather_id FROM race_forms WHERE id=$1", summary[idx].RaceFormID).Scan(&raceform)
		if err != nil {
			return []repository.RaceSummary{}, fmt.Errorf("unable to fetch raceform with id %s failed: %w", summary[idx].RaceFormID, err)
		}

		// Distance type
		var distanceType struct {
			ID        string `db:"id"`
			Name      string `db:"name"`
			ShortName string `db:"short_name"`
		}
		err = p.pool.QueryRow(context.Background(), "SELECT id, name, short_name FROM distance_types WHERE id=$1", raceform.DistanceTypeID).Scan(&distanceType)
		if err != nil {
			return []repository.RaceSummary{}, fmt.Errorf("unable to fetch distance_type with id %s failed: %w", raceform.DistanceTypeID, err)
		}

		// Track condition
		var trackCondition struct {
			ID        string `db:"id"`
			Name      string `db:"name"`
			ShortName string `db:"short_name"`
		}
		err = p.pool.QueryRow(context.Background(), "SELECT id, name, short_name FROM track_conditions WHERE id=$1", raceform.TrackConditionID).Scan(&trackCondition)
		if err != nil {
			return []repository.RaceSummary{}, fmt.Errorf("unable to fetch track_conditions with id %s failed: %w", raceform.TrackConditionID, err)
		}

		// Weather
		var weather struct {
			IconURI   string `db:"icon_uri"`
			ID        string `db:"id"`
			Name      string `db:"name"`
			ShortName string `db:"short_name"`
		}
		err = p.pool.QueryRow(context.Background(), "SELECT icon_uri, id, name, short_name FROM weather_conditions WHERE id=$1", raceform.WeatherID).Scan(&weather)
		if err != nil {
			return []repository.RaceSummary{}, fmt.Errorf("unable to fetch weather with id %s failed: %w", raceform.WeatherID, err)
		}
		// Put all the data into the DTO
		rs := repository.RaceSummary{}
		rs.ID = summary[idx].ID
		rs.AdvertisedStart.Seconds = summary[idx].Seconds
		rs.CategoryID = summary[idx].CategoryID
		rs.MeetingID = summary[idx].MeetingID
		rs.MeetingName = summary[idx].MeetingName
		rs.RaceFormID = summary[idx].RaceFormID
		rs.RaceID = summary[idx].RaceID
		rs.RaceName = summary[idx].RaceName
		rs.RaceNumber = summary[idx].RaceNumber
		rs.VenueCountry = summary[idx].VenueCountry
		rs.VenueID = summary[idx].VenueID
		rs.VenueName = summary[idx].VenueName
		rs.VenueState = summary[idx].VenueState
		rs.RaceForm.AdditionalData = raceform.AdditionalData
		rs.RaceForm.Distance = raceform.Distance
		rs.RaceForm.DistanceTypeID = raceform.DistanceTypeID
		rs.RaceForm.Generated = raceform.Generated
		rs.RaceForm.RaceComment = raceform.RaceComment
		rs.RaceForm.RaceCommentAlternative = raceform.RaceCommentAlternative
		rs.RaceForm.SilkBaseURL = raceform.SilkBaseURL
		rs.RaceForm.TrackConditionID = raceform.TrackConditionID
		rs.RaceForm.WeatherID = raceform.WeatherID
		rs.RaceForm.DistanceType.ID = distanceType.ID
		rs.RaceForm.DistanceType.Name = distanceType.Name
		rs.RaceForm.DistanceType.ShortName = distanceType.ShortName
		rs.RaceForm.TrackCondition.ID = trackCondition.ID
		rs.RaceForm.TrackCondition.Name = trackCondition.Name
		rs.RaceForm.TrackCondition.ShortName = trackCondition.ShortName
		rs.RaceForm.Weather.ID = weather.ID
		rs.RaceForm.Weather.IconURI = weather.IconURI
		rs.RaceForm.Weather.Name = weather.Name
		rs.RaceForm.Weather.ShortName = weather.ShortName
		results = append(results, rs)
	}
	log.Printf("Postgres returning results %v", results)
	return results, nil
}
