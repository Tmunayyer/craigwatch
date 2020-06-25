package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	// postgres driver
	"github.com/lib/pq"
)

// TODO: Cleanup this interface a bit, start centing things around table
// TODO: rename montior/URL stuff to search
type connection interface {
	connect() error
	shutdown() error
	testConnection() error
	applySchema() error
	saveSearch(clSearch) (clSearch, error)
	deleteSearch(id int) error
	getAllSearches() ([]clSearch, error)
	saveListings(monitorID int, listings []clListing) error
	deleteListings(monitorID int) error
	getListings(id int) ([]clListing, error)
}

type client struct {
	db               *sql.DB
	connectionString string
}

// newDBClient will return a connected client with reference to
// the connection string. The connection string is constructed using
// environment variables.
func newDBClient() connection {
	c := client{}
	err := c.connect()
	if err != nil {
		panic(err)
	}
	err = c.applySchema()
	if err != nil {
		panic(err)
	}

	return &c
}

// Connect to the database
// connect will set up the connection and store the connection string
// for reference.
func (c *client) connect() error {
	var connectionString = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"),
	)

	connection, err := sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("connection to pg failed: %v", err)
	}

	c.db = connection
	c.connectionString = connectionString

	fmt.Println("postgres connection established...")
	return nil
}

// Shutdown to close the connection to PostgreSQL
func (c *client) shutdown() error {
	err := c.db.Close()
	if err != nil {
		return fmt.Errorf("error closing connection to pg: %v", err)
	}
	fmt.Println("closing connection to postgres...")
	return nil
}

// TestConnection performs a ping query on PostgreSQL
func (c *client) testConnection() error {
	return c.db.Ping()
}

func (c *client) applySchema() error {
	file, err := ioutil.ReadFile("./schema.sql")
	if err != nil {
		return fmt.Errorf("error reading schema file: %v", err)
	}

	_, err = c.db.Query(string(file))
	if err != nil {
		return fmt.Errorf("schema query error: %v", err)
	}

	return nil
}

// =======================
// ===== Models
// =======================

type clSearch struct {
	ID        int
	Name      string
	URL       string
	Confirmed bool
	Interval  int
	CreatedOn time.Time
}

type clListing struct {
	ID           int
	MonitorID    int
	DataPID      string
	DataRepostOf string
	Date         time.Time
	Title        string
	Link         string
	Price        int
	Hood         string
}

// =======================
// ===== Queries
// =======================

func (c *client) saveSearch(data clSearch) (clSearch, error) {
	output := clSearch{}

	rows, err := c.db.Query(`
		insert into monitor
			(name, url, confirmed, created_on)
		values
			($1, $2, $3, Now())
		returning *
	`, data.Name, data.URL, false)
	defer rows.Close()

	if err != nil {
		return output, err
	}

	for rows.Next() {
		err := rows.Scan(
			&output.ID,
			&output.Name,
			&output.URL,
			&output.Confirmed,
			&output.Interval,
			&output.CreatedOn,
		)
		if err != nil {
			return output, err
		}
	}

	err = rows.Err()
	if err != nil {
		return output, err
	}

	return output, nil
}

func (c *client) getAllSearches() ([]clSearch, error) {
	output := []clSearch{}

	rows, err := c.db.Query(`
		select
			id,
			name,
			url,
			confirmed,
			interval,
			created_on
		from 
			monitor
	`)

	if err != nil {
		return output, err
	}

	for rows.Next() {
		q := clSearch{}
		err := rows.Scan(
			&q.ID,
			&q.Name,
			&q.URL,
			&q.Confirmed,
			&q.Interval,
			&q.CreatedOn,
		)
		if err != nil {
			return output, err
		}

		output = append(output, q)
	}

	err = rows.Err()
	if err != nil {
		return output, err
	}

	return output, nil
}

func (c *client) deleteSearch(id int) error {
	_, err := c.db.Query(`
		delete from 
			monitor
		where
			id = $1
	`, id)

	if err != nil {
		return err
	}

	return nil
}

func (c *client) saveListings(monitorID int, listings []clListing) error {
	txn, err := c.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn("listing", "monitor_id", "data_pid", "data_repost_of", "date", "title", "link", "price", "hood"))
	if err != nil {
		fmt.Println("from the formatting")
		return err
	}

	for _, l := range listings {
		_, err = stmt.Exec(monitorID, l.DataPID, l.DataRepostOf, l.Date, l.Title, l.Link, l.Price, l.Hood)

		if err != nil {
			fmt.Println("from the executing")
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *client) deleteListings(monitorID int) error {
	_, err := c.db.Query(`
		delete from 
			listing
		where
			monitor_id = $1
	`, monitorID)

	if err != nil {
		return err
	}

	return nil
}

func (c *client) getListings(monitorID int) ([]clListing, error) {
	output := []clListing{}

	rows, err := c.db.Query(`
		select
			*
		from 
			listing
		where
			monitor_id = $1
		order by
			date desc;
	`, monitorID)

	if err != nil {
		return output, err
	}

	for rows.Next() {
		q := clListing{}
		err := rows.Scan(
			&q.ID,
			&q.MonitorID,
			&q.DataPID,
			&q.DataRepostOf,
			&q.Date,
			&q.Title,
			&q.Link,
			&q.Price,
			&q.Hood,
		)
		if err != nil {
			return output, err
		}

		output = append(output, q)
	}

	err = rows.Err()
	if err != nil {
		return output, err
	}

	return output, nil
}
