CREATE USER ladbrokes;
CREATE DATABASE race_data_db owner ladbrokes;
GRANT ALL PRIVILEGES ON DATABASE race_data_db TO ladbrokes;
\connect race_data_db ladbrokes
CREATE TABLE IF NOT EXISTS race_summaries
(
	id TEXT PRIMARY KEY,
	advertised_start int,
	category_id TEXT,
	meeting_id TEXT,
	meeting_name TEXT,
	race_form int,
	race_id TEXT,
	race_name TEXT,
	race_number int,
	venue_country TEXT,
	venue_id TEXT,
	venue_name TEXT,
	venue_state TEXT
);
	/* title TEXT, */
	/* pub_date DATE, */
	/* body TEXT, */
	/* tags TEXT[] */
CREATE TABLE IF NOT EXISTS race_forms
(
	id SERIAL PRIMARY KEY,
	distance int
	distance_type_id TEXT,
	generated TEXT,
	race_comment TEXT,
	race_comment_alternative TEXT,
	silk_base_url TEXT,
	track_condition TEXT,
	weather_id TEXT
);
CREATE TABLE IF NOT EXISTS distance_types
(
	id TEXT PRIMARY KEY,
	name TEXT,
	short_name TEXT
);
CREATE TABLE IF NOT EXISTS track_conditions
(
	id TEXT PRIMARY KEY,
	name TEXT,
	short_name TEXT
);
CREATE TABLE IF NOT EXISTS weather_conditions
(
	id TEXT PRIMARY KEY,
	icon_uri TEXT,
	name TEXT,
	short_name TEXT
);
ALTER USER ladbrokes WITH PASSWORD 'hu8jmn3';
