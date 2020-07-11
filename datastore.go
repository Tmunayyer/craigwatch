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

// TODO: rename monitor table to search
type connection interface {
	connect() error
	shutdown() error
	testConnection() error
	applySchema() error

	saveSearch(clSearch) (clSearch, error)
	deleteSearch(id int) error
	getSearch(searchID int) (clSearch, error)
	getSearchMulti() ([]clSearch, error)

	saveListingMulti(monitorID int, listings []clListing) error
	deleteListingMulti(monitorID int) error
	getListingMulti(id int) ([]clListing, error)
	getListingMultiAfter(id int, unixDate int64) ([]clListing, error)

	getSearchActivity(searchID int) (searchActivity, error)
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

const (
	duplicateURLErrorMessage = `pq: duplicate key value violates unique constraint "search_url_key"`
)

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
	ID             int
	Name           string
	URL            string
	Confirmed      bool
	Interval       int
	CreatedOn      time.Time
	UnixCutoffDate int
}

type searchActivity struct {
	ID        int
	InSeconds int
	InMinutes float32
	InHours   float32
}

type clListing struct {
	ID           int
	SearchID     int
	DataPID      string
	DataRepostOf string
	UnixDate     int64
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
	insert into search
		(name, url, created_on)
	values
		($1, $2, Now())
	returning *
	`, data.Name, data.URL)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&output.ID,
			&output.Name,
			&output.URL,
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

func (c *client) getSearch(searchID int) (clSearch, error) {
	output := clSearch{}

	rows, err := c.db.Query(`
		select
			s.*,
			coalesce(l.unix_cutoff_date, '0')
		from
			search s
		left join
			(
				select
					search_id,
					max(unix_date) as "unix_cutoff_date"
				from listing
				group by search_id
			) l
		on l.search_id = s.id
		where s.id = $1
	`, searchID)

	if err != nil {
		return output, err
	}

	for rows.Next() {
		err := rows.Scan(
			&output.ID,
			&output.Name,
			&output.URL,
			&output.CreatedOn,
			&output.UnixCutoffDate,
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

func (c *client) getSearchMulti() ([]clSearch, error) {
	output := []clSearch{}

	rows, err := c.db.Query(`
		select
			s.*,
			coalesce(l.unix_cutoff_date, '0')
		from
			search s
		left join
			(
				select
					search_id,
					max(unix_date) as "unix_cutoff_date"
				from listing
				group by search_id
			) l
		on l.search_id = s.id
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
			&q.CreatedOn,
			&q.UnixCutoffDate,
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
			search
		where
			id = $1
	`, id)

	if err != nil {
		return err
	}

	return nil
}

func (c *client) saveListingMulti(searchID int, listings []clListing) error {
	txn, err := c.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn("listing", "search_id", "data_pid", "data_repost_of", "unix_date", "title", "link", "price", "hood"))
	if err != nil {
		return err
	}

	for _, l := range listings {
		_, err = stmt.Exec(searchID, l.DataPID, l.DataRepostOf, l.UnixDate, l.Title, l.Link, l.Price, l.Hood)

		if err != nil {
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

func (c *client) deleteListingMulti(searchID int) error {
	_, err := c.db.Query(`
		delete from 
			listing
		where
			search_id = $1
	`, searchID)

	if err != nil {
		return err
	}

	return nil
}

func (c *client) getListingMulti(searchID int) ([]clListing, error) {
	output := []clListing{}

	rows, err := c.db.Query(`
		select
			*
		from 
			listing
		where
			search_id = $1
		order by
			unix_date desc;
	`, searchID)

	if err != nil {
		return output, err
	}

	for rows.Next() {
		q := clListing{}
		err := rows.Scan(
			&q.ID,
			&q.SearchID,
			&q.DataPID,
			&q.DataRepostOf,
			&q.UnixDate,
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

func (c *client) getListingMultiAfter(searchID int, unixDate int64) ([]clListing, error) {
	output := []clListing{}

	rows, err := c.db.Query(`
		select
			id,
			search_id,
			data_pid,
			data_repost_of,
			unix_date,
			title,
			link,
			price,
			hood
		from 
			listing
		where
			search_id = $1
		and
			unix_date > $2
		order by
			unix_date desc;
	`, searchID, unixDate)

	if err != nil {
		return output, err
	}

	for rows.Next() {
		q := clListing{}
		err := rows.Scan(
			&q.ID,
			&q.SearchID,
			&q.DataPID,
			&q.DataRepostOf,
			&q.UnixDate,
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

func (c *client) getSearchActivity(searchID int) (searchActivity, error) {
	output := searchActivity{}

	rows, err := c.db.Query(`
		with order_by_date as (
			select 
				l.*
			from listing l
			left join search s
			on l.search_id = s.id
			where s.id = $1
			order by unix_date desc
		),
		
		time_between_posts as ( 
			select
				ao.search_id,
				unix_date - lead(unix_date, 1) over (order by unix_date desc) as time_between
			from order_by_date ao
		)
		
		select
			tbp.search_id,
			round(avg(tbp.time_between)) as in_seconds,
			round(avg(tbp.time_between) / 60, 2) as in_minutes,
			round(avg(tbp.time_between) / 60 / 60, 2) as in_hours
		from time_between_posts tbp
		group by tbp.search_id;
	`, searchID)

	if err != nil {
		return output, err
	}

	for rows.Next() {
		err := rows.Scan(
			&output.ID,
			&output.InSeconds,
			&output.InMinutes,
			&output.InHours,
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
