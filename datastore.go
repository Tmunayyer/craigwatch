package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	// postgres driver
	_ "github.com/lib/pq" // here
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
	getSearchActivityByHour(searchID int) ([]activityByHour, error)
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
	var connection *sql.DB
	var err error
	if os.Getenv("MODE") == "development" {
		var connectionString = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("PGHOST"),
			os.Getenv("PGPORT"),
			os.Getenv("PGUSER"),
			os.Getenv("PGPASSWORD"),
			os.Getenv("PGDATABASE"),
		)
		c.connectionString = connectionString
		connection, err = sql.Open("postgres", connectionString)
	} else if os.Getenv("MODE") == "production" {
		connection, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	}
	if err != nil {
		return fmt.Errorf("connection to pg failed: %v", err)
	}

	c.db = connection

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
	UnixCutoffDate int64
	Timezone       string
	TotalListings  int
}

type searchActivity struct {
	ID              int
	InSeconds       int
	InMinutes       float32
	InHours         float32
	RepostInSeconds int
	RepostInMinutes float32
	RepostInHours   float32
	TotalCount      int
	RepostCount     int
	PercentRepost   float32
}

type activityByHour struct {
	Label       time.Time
	TopUnixDate int64
	BotUnixDate int64
	Count       int
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

	stmt := `
		insert into search
			(name, url, created_on, timezone)
		values
			($1, $2, now(), $3)
		returning *
	`
	rows, err := c.db.Query(stmt, data.Name, data.URL, data.Timezone)
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
			&output.Timezone,
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

	stmt := `
		select
			s.*,
			coalesce(l.unix_cutoff_date, '0'),
			coalesce(l.total_listings, 0)
		from
			search s
		left join
			(
				select
					search_id,
					max(unix_date) as "unix_cutoff_date",
					count(*) as "total_listings"
				from listing
				group by search_id
			) l
		on l.search_id = s.id
		where s.id = $1
	`
	rows, err := c.db.Query(stmt, searchID)
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
			&output.Timezone,
			&output.UnixCutoffDate,
			&output.TotalListings,
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

	stmt := `
		select
			s.*,
			coalesce(l.unix_cutoff_date, '0'),
			coalesce(l.total_listings, 0)
		from
			search s
		left join
			(
				select
					search_id,
					max(unix_date) as "unix_cutoff_date",
					count(*) as "total_listings"
				from listing
				group by search_id
			) l
		on l.search_id = s.id
	`
	rows, err := c.db.Query(stmt)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		q := clSearch{}
		err := rows.Scan(
			&q.ID,
			&q.Name,
			&q.URL,
			&q.CreatedOn,
			&q.Timezone,
			&q.UnixCutoffDate,
			&q.TotalListings,
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
	stmt := `
		delete from 
			search
		where
			id = $1
	`
	rows, err := c.db.Query(stmt, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (c *client) saveListingMulti(searchID int, listings []clListing) error {
	insertStatement := `
		insert into 
			listing 
			(search_id, data_pid, data_repost_of, unix_date, title, link, price, hood) 
		values
	`
	conflictStatement := `
		on conflict (data_pid)
		do update set
			data_repost_of = excluded.data_repost_of,
			unix_date = excluded.unix_date,
			title = excluded.title,
			link = excluded.link,
			price = excluded.price,
			hood = excluded.hood;
	`

	valueStatement := ""
	values := make([]interface{}, 0, len(listings)*8+1)
	vIndex := 1
	for i, listing := range listings {

		row := ""

		row += "$" + strconv.Itoa(vIndex) + ","
		values = append(values, searchID)
		vIndex++

		row += "$" + strconv.Itoa(vIndex) + ","
		values = append(values, listing.DataPID)
		vIndex++

		row += "$" + strconv.Itoa(vIndex) + ","
		values = append(values, listing.DataRepostOf)
		vIndex++

		row += "$" + strconv.Itoa(vIndex) + ","
		values = append(values, listing.UnixDate)
		vIndex++

		row += "$" + strconv.Itoa(vIndex) + ","
		values = append(values, listing.Title)
		vIndex++

		row += "$" + strconv.Itoa(vIndex) + ","
		values = append(values, listing.Link)
		vIndex++

		row += "$" + strconv.Itoa(vIndex) + ","
		values = append(values, listing.Price)
		vIndex++

		row += "$" + strconv.Itoa(vIndex)
		values = append(values, listing.Hood)
		vIndex++

		row = "(" + row + ")"

		if i < len(listings)-1 {
			row += ","
		}

		valueStatement += row

	}

	statement := insertStatement + valueStatement + conflictStatement

	rows, err := c.db.Query(statement, values...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (c *client) deleteListingMulti(searchID int) error {
	rows, err := c.db.Query(`
		delete from 
			listing
		where
			search_id = $1
	`, searchID)
	defer rows.Close()

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
	defer rows.Close()

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
	defer rows.Close()

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
		with total_listings as (
			select 
				l.*
			from listing l
			left join search s
			on l.search_id = s.id
			where s.id = $1
			order by unix_date desc
		),
		
		repost_listings as (
			select
				*
			from total_listings
			where total_listings.data_repost_of <> ''
		),
		
		time_between_posts as ( 
			select
				tl.search_id,
				unix_date - lead(unix_date, 1) over (order by unix_date desc) as time_between
			from total_listings tl
		),
		
		time_between_reposts as ( 
			select
				rpl.search_id,
				unix_date - lead(unix_date, 1) over (order by unix_date desc) as time_between
			from repost_listings rpl
		),
		
		post_frequency as (
			select
				tbp.search_id,
				count(*) as total_count,
				round(avg(tbp.time_between)) as in_seconds,
				round(avg(tbp.time_between) / 60, 2) as in_minutes,
				round(avg(tbp.time_between) / 60 / 60, 2) as in_hours
			from time_between_posts tbp
			group by tbp.search_id
		),
		
		repost_frequency as (
			select
				tbrp.search_id,
				count(*) as repost_count,
				round(avg(tbrp.time_between)) as in_seconds,
				round(avg(tbrp.time_between) / 60, 2) as in_minutes,
				round(avg(tbrp.time_between) / 60 / 60, 2) as in_hours
			from time_between_reposts tbrp
			group by tbrp.search_id
		)
		
		select
			pf.search_id,
			pf.in_seconds,
			pf.in_minutes,
			pf.in_hours,
			rpf.in_seconds as repost_in_seconds,
			rpf.in_minutes as repost_in_minutes,
			rpf.in_hours as repost_in_hours,
			pf.total_count,
			rpf.repost_count,
			round(rpf.repost_count::numeric / pf.total_count::numeric, 2) as percent_reposts
		from post_frequency as pf
		full join repost_frequency rpf on pf.search_id = rpf.search_id
	`, searchID)
	defer rows.Close()

	if err != nil {
		return output, err
	}

	for rows.Next() {
		err := rows.Scan(
			&output.ID,
			&output.InSeconds,
			&output.InMinutes,
			&output.InHours,
			&output.RepostInSeconds,
			&output.RepostInMinutes,
			&output.RepostInHours,
			&output.TotalCount,
			&output.RepostCount,
			&output.PercentRepost,
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

func (c *client) getSearchActivityByHour(searchID int) ([]activityByHour, error) {
	output := []activityByHour{}

	rows, err := c.db.Query(`
		with total_listings as (
			select 
				l.*
			from listing l
			left join search s
			on l.search_id = s.id
			where s.id = $1
			order by unix_date desc
		),
		
		max_date as (
			select 
				unix_date as maximum,
				EXTRACT(EPOCH FROM NOW())::bigint as "now"
			from total_listings
			limit 1
		),
		
		cutoff_dates as (
			select
				case 
					when tl.unix_date = (select "now" from max_date)
					then to_timestamp(tl.unix_date) 
					else to_timestamp((select "now" from max_date) - ((row_number() over (order by tl.unix_date desc) - 1) * 3600))
				end "label",
				case
					when tl.unix_date = (select "now" from max_date)
					then tl.unix_date
					else (select "now" from max_date) - ((row_number() over (order by tl.unix_date desc) - 1) * 3600)
				end "top_unix_date",
				(select "now" from max_date) - ((row_number() over (order by tl.unix_date desc)) * 3600) as "bot_unix_date"
			from total_listings tl
			join (select search_id, unix_date as "max" from total_listings limit 1) mu on tl.search_id = mu.search_id
			limit 48
		)
		
		select
			cd.label,
			cd.top_unix_date,
			cd.bot_unix_date,
			(select count(*) filter (where tl.unix_date < cd.top_unix_date - 1 and tl.unix_date > cd.bot_unix_date - 1) from total_listings tl)
		from cutoff_dates cd
	`, searchID)
	defer rows.Close()

	if err != nil {
		return output, err
	}

	for rows.Next() {
		set := activityByHour{}

		err := rows.Scan(
			&set.Label,
			&set.TopUnixDate,
			&set.BotUnixDate,
			&set.Count,
		)
		if err != nil {
			return output, err
		}

		output = append(output, set)
	}

	err = rows.Err()
	if err != nil {
		return output, err
	}

	return output, nil
}
