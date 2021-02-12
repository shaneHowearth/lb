CREATE USER ladbrokes;
CREATE DATABASE race_data_db owner ladbrokes;
GRANT ALL PRIVILEGES ON DATABASE race_data_db TO ladbrokes;
\connect race_data_db ladbrokes
CREATE TABLE IF NOT EXISTS race_summaries
(
	id SERIAL PRIMARY KEY,
	category_id UUID,
	meeting_id,
	meeting_name
	race_form
	race_id
	race_name
	race_number
	venue_country
	venue_id
	venue_name
	venue_state
);
	/* title TEXT, */
	/* pub_date DATE, */
	/* body TEXT, */
	/* tags TEXT[] */
CREATE TABLE IF NOT EXISTS race_forms
(
	id SERIAL PRIMARY KEY,
	distance
	distance_type
	generated
	race_comment TEXT
	race_comment_alternative
	silk_base_url
	track_condition
	weather_id
);
CREATE TABLE IF NOT EXISTS distance_types
(
	id SERIAL PRIMARY KEY,
	name
	short_name
);
CREATE TABLE IF NOT EXISTS track_conditions
(
	id SERIAL PRIMARY KEY,
	name
	short_name
);
CREATE TABLE IF NOT EXISTS weather_conditions
(
	id SERIAL PRIMARY KEY,
	icon_uri
	name
	short_name
);
ALTER USER ladbrokes WITH PASSWORD 'hu8jmn3';
