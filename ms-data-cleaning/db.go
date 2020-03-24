package main

import (
	"database/sql"
	"log"
)

//Database controls database functionality
type Database struct {
	db *sql.DB
}

//Row represents a row in the database for line counts
type Row struct {
	Key   string
	Value int
}

//StoreValue stores a key and value in a database
func (d *Database) StoreValue(key string, value int) error {

	sqlStmt := `
	INSERT INTO SampleTable (key, value)
	VALUES (?,?)
	`
	statement, err := d.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = statement.Exec(key, value)
	if err != nil {
		log.Println("Failed to execute sql", err)
		return err
	}

	return nil

}

//fetchData allows you to fetch data from db.
func (d *Database) fetchValues(fname string) ([]Row, error) {
	rows, err := d.db.Query("SELECT * FROM SampleTable ")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rs := []Row{}
	for rows.Next() {
		r := Row{}
		err := rows.Scan(&r.Key,
			&r.Value)
		if err != nil {
			log.Println("Failed to fetch data from db: ", err)
			return rs, err
		}
		rs = append(rs, r)
	}
	return rs, nil
}

//dbinit function will create a table for use for this microservice.
//Change this to include the table you need for your service
func (d *Database) dbInit() {

	//create browser table
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS SampleTable (
		key TEXT,
		VALUE int
		)
	`

	statement, _ := d.db.Prepare(sqlStmt)
	statement.Exec()
}
