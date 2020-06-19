# craigwatch
A tool to monitor craigslist searches and so much more.

## Table of Contents

- [Setting Up](#setting-up)

## Setting Up

1. [Install PostgreSQL](https://www.postgresql.org/download/)
    - Connect to PostgreSQL shell

    ```
    $ psql postgres
    ```

    - Create the database

    ```
    $ create database craigwatch;
    ```

2. [Install Go](https://golang.org/doc/install)
3. Navigate to the repo's root
4. Start the server

    ```
    go run .
    ```
5. Navigate to localhost:3000 in your browser