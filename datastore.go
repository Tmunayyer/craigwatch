package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

	// postgres driver
	_ "github.com/lib/pq"
)

type connection interface {
	connect() error
	shutdown() error
	testConnection() error
	applySchema() error
	saveURL(craigslistQuery) (craigslistQuery, error)
}

type client struct {
	db               *sql.DB
	connectionString string
}

// newDBClient will return a connected client with reference to
// the connection string. The connection string is constructed using
// environment variables.
func newDBClient() (connection, error) {
	c := client{}
	err := c.connect()
	if err != nil {
		return nil, fmt.Errorf("error connecting to pg: %v", err)
	}
	err = c.applySchema()
	if err != nil {
		return nil, fmt.Errorf("error applying schema: %v", err)
	}

	return &c, nil
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

type craigslistQuery struct {
	ID        int
	Email     string
	URL       string
	Confirmed bool
	Interval  int
	CreatedOn string
	PolledOn  string
}

// =======================
// ===== Queries
// =======================

func (c *client) saveURL(data craigslistQuery) (craigslistQuery, error) {
	output := craigslistQuery{}

	rows, err := c.db.Query(`
		insert into monitor
			(email, url, confirmed)
		values
			($1, $2, $3)
		returning *
	`, data.Email, data.URL, false)
	defer rows.Close()

	if err != nil {
		return output, err
	}

	for rows.Next() {
		err := rows.Scan(&output)
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
