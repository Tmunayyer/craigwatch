# craigwatch
A tool to monitor craigslist searches and so much more.

## Table of Contents

- [The Application](https://craigwatch.herokuapp.com/) note: may take a few seconds to wake up
- [About](#about)
- [Setting Up](#setting-up)
- [API Documentation](./documentation/api/table_of_contents.MD)
- [Takeaways](#takeaways)

## About

This goal of this project was to create something interesting while learning new technologies. Prior to this project my experience in Go was limited to a library I wrote called [go-craigslist](https://github.com/Tmunayyer/go-craigslist) which crawls and parses craigslist.org. 

I wanted a bit more experience centered around a web application. Since I already had a tool to do much of the heavy lifting, I thought it might be interesting to analyze the craigslist postings for a given search. The idea resulted in this project.

I also decided to try something new on the frontend. I had only really worked in React but always heard about Vue and how great it was. I decided to use this opportunity to give it try and implement it here.

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

## Takeaways

### Frontend

The frontend of the application was really a whole new approach to building user interfaces. The overall structure of the application isn't too different than what it would have been since both React and Vue are component centric frameworks. What is different was the approach to building a specific component.

When learning about React, something I remember in detail was that JSX, the standard “templating” language of React provides the full power of Javascript. This was notably missing from Vue. Vue’s templates keep things cleaner overall but without the knowledge of the syntax, the time it took to develop was very slow at first. This was really the starkest contrast between the two however given sufficient time with both, im sure would be a non-factor.

### Backend

This was the first web server I have built using Go. I did not use a framework as I wanted to get as comfortable with the standard http module as possible. Really, the extent of the server is just relaying data and there is no real significant logic required.

Because I did not use a framework, it's a bit hard to compare to Node/Express. One thing that really stands out is the testing capabilities of Go. Having a standard testing library and the support for it gives it a huge leg up over anything I have found in Javascript so far.

Another standout for me was the introduction of types. Earlier in my career I would come across instances of functions returning undefined causing property lookup errors down the line. This simply doesn't happen when the compiler enforces function signatures. What started as a frustrating and new process of battling with the compiler, I quickly realized that most silly errors common in javascript don't exist at runtime with Go.


### Additional

I really made an effort to set up all environments myself. Webpack is a tool I have used my entire career and I feel like I learned more about this project than the past year just by writing a silly custom plugin. 

The webpack config file contains a class called TinyTidyPlugin which was the result of battling browser caching and bundle naming. Basically, the browser was caching my bundle.js file which was inhibiting the development experience while using webpack’s watch feature. To stop this from happening, webpack supports adding a hash to the bundle file name. Now that it's producing a new bundle file each time, they would pile up in the directory.

After noticing this and letting it bother me for a few minutes I figured it would be an easy “pass this function to delete the old bundle”. While it wasn't quite that simple, I found writing the plugin to be a small but extremely rewarding task. I learned quite a bit about the webpack process and lifecycle hooks. Other plugins are no longer a black box to me and, most importantly, the old bundles are deleted when a new one is produced.

