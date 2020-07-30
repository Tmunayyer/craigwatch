# craigwatch
A tool to monitor craigslist searches and so much more.

## Table of Contents

- [The Application](#https://craigwatch.herokuapp.com/) note: may take a few seconds to wake up
- [Setting Up](#setting-up)
- [API Documentation](./documentation/api/table_of_contents.MD)

## Setting Up

1. Install PostgreSQL ([official docs](https://www.postgresql.org/download/))
    - Connect to PostgreSQL shell

    ```
    $ psql postgres
    ```

    - Create the database

    ```
    $ create database craigwatch;
    ```

2. Install Go ([official docs](https://golang.org/doc/install))
3. Set up your environment variables in .env file
    ```
    MODE=development

    PGHOST=localhost
    PGPORT=<<your desired port>>
    PGUSER=<<your username>>
    PGPASSWORD=<<your password>>
    PGDATABASE=craigwatch
    ```
4. Navigate to the repo's root in the terminal
5. Start the server

    ```
    npm run server-dev
    ```
6. Start webpack
    ```
    npm run build-dev
    ```
7. Navigate to localhost:3000 in your browser