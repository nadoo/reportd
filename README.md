reportd
=====

An easy report server for generate html tables from database.

## Install

	go get -u github.com/nadoo/reportd

## Build

	cd $GOPATH/src/github.com/nadoo/reportd
	go build

## Config

	# Listen Address: "ip:port" OR ":port"(listen on all interfaces)
	Listen = ":8080"

	# Database Type: "mysql" OR "postgres"
	DBType = "mysql"

	# Database connection string
	DBConn = "user:password@tcp(127.0.0.1:3306)/database?charset=utf8&autocommit=true" # mysql
	#DBConn = "postgres://user:password@127.0.0.1:5432/database?sslmode=disable" # postgres

	[[Reports]]
	Title = "Report 1"
	Sql = """SELECT
			s.song_id AS SongID
			,s.song_name
		FROM m_song s
		LIMIT 1;
	"""

	[[Reports]]
	Title = "Report 2"
	Params = true
	Sql = """SELECT
				s.song_id,
				s.song_name
			FROM m_song s
			WHERE s.song_artist_only IN (:artists);
	"""

	#[[Reports]]
	#Title = "Report 3"
	#Params = true
	#Sql = """SELECT
	#			s.song_id,
	#			s.song_name
	#		FROM m_song s
	#		WHERE s.song_artist_only = :artist1 OR s.song_artist_only = :artist2;
	#"""


## Usage
1. Run the program:
> nohup ./reportd &

2. View in the browser: 
	- html output
		> http://127.0.0.1:8080
	- html output with query parameters
		> http://127.0.0.1:8080?artists=she&artists=her
	- json output
		> http://127.0.0.1:8080/json
	- json output with query parameters
		> http://127.0.0.1:8080/json?artists=she&artists=her