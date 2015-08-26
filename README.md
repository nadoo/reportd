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
	Sql = "SELECT s.song_id AS SongID, s.song_name FROM m_song s LIMIT 1;"

	[[Reports]]
	Title = "Report 2"
	Sql = "SELECT s.song_id, s.song_name, s.song_artist_only FROM m_song s LIMIT 5;"

## Usage
1. Run the program:
> nohup ./reportd &

2. View in the browser: http://127.0.0.1:8080
