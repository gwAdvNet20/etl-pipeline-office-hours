package main

import (
	"database/sql"
	"log"
)

//Database controls database functionality
type Database struct {
	db *sql.DB
}

//LineCountRow represents a row in the database for line counts
type LineCountRow struct {
	Key   string
	Count int
}

//Store stores a logfile in a database
func (d *Database) StoreLogLine(lf LogFile) {

	sqlStmt := `
	INSERT INTO logs (name, raw_log, 
		remote_addr,
		time_local,
		request_type,
		request_path,
		status,
		body_bytes_sent,
		http_referer,
		http_user_agent,
		created
		) VALUES (?,?, ?,?,?,?,?,?,?,?,?)
	`
	statement, err := d.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range lf.Logs {

		statement.Exec(v.Name, v.RawLog, v.RemoteAddr, v.TimeLocal, v.RequestType, v.RequestPath, v.Status, v.BodyBytesSent, v.HTTPReferer, v.HTTPUserAgent, v.Created)
	}
}

func (d *Database) dbInit() {

	//create logs table
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS logs (
		name TEXT NOT NULL,
		raw_log TEXT NOT NULL UNIQUE,
		remote_addr TEXT,
		time_local TEXT,
		request_type TEXT,
		request_path TEXT,
		status INTEGER,
		body_bytes_sent INTEGER,
		http_referer TEXT,
		http_user_agent TEXT,
		created DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`

	statement, _ := d.db.Prepare(sqlStmt)
	statement.Exec()
}
