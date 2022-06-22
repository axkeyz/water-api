# My water is down again!

This is a silly project that hopes to prove some regions in Auckland get more water outages than others (I swear I don't live that rurally!)

Unfortunately there is no public API for previous water outages, so the data collected by this is largely incomplete.

But one day, I'll prove it! Maybe.

## APIs

1. Main API, available at the root of the server (/).

    It comes with the following (query) parameters - these narrow down data.
    - outage_type
    - before_start_date & after_start_date
    - before_end_date & after_end_date
    - suburb
    - street
    - location (needs longitude + latitude + radius)

    *Example 1*: /?outage_type=Planned&suburb=Remuera 
    Returns results of all planned outages in Remuera.

    *Example 2*: /?location=true&longitude=174.762415&latitude=-36.855109&radius=2000 
    Returns all outages that happened within 2 km (2000 m) of Queen Street (174l762416, -36.855109).

2. Count API, available at /count.

    Same query parameters as above (narrow down results).

    It also comes with chainable "get" parameters which divide counts by those categories. These are the same as the query parameters, with an extra total_hours value.

    *Example 1*: /count?get=outage_type 
    Counts up the unplanned and planned outages.

    *Example 2*: /count?get=suburb&outage_type=Unplanned&get=total_hours
    Gets a count of all outages per suburb that are unplanned. It also gets the total hours.

### Pagination

Comes with "limit" & "offset" parameters, where limit is the total number of items returned and offset is the number of items to skip before counting the needed data.

It *needs* a "sort" parameter.

*Example*: /count?get=total_hours&get=suburb&sort=total_outages%20desc&sort=suburb&limit=10&offset=10
Gets total outages & hours of 10 suburbs, descending sorted by total number of outages (unluckiest first). Only 10 suburbs are returned, ranking 11-20 of the most unluckiest.

## Installation instructions

1. Copy the following files & make changes as needed:
    - docker-compose.yml
    - .env-example: Rename to .env when done
        - SRC_API: Original outage API (replace for testing purposes)
    - docker_postgres_init.sql
2. Pull prepared image from DockerHub and start: ```docker-compose up -d```
3. Navigate to localhost:APP_PORT (whatever you set up in the .env file)

## Data collection

Data is collected every 1 hour.

## Live version

There's one at: https://water.aileenhuang.dev/