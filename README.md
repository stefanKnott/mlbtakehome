# Take Home Assignment

## The API
`/api/v1/schedule?teamId=<id>&date=<YYYY-MM-DD>`

This API allows for a client to query the backend for scheduled games for a given date.  The response payload is ordered such that the games for the requested team with the corresponding `teamId` are listed first.

### Query Parameters
* `teamId`: an integer value for a valid MLB team (ie. 141).  A list of valid teams for the 2024 season can be found [here](https://statsapi.mlb.com/api/v1/teams?season=2024&sportId=1).
* `date`: a string value of the format `YYYY-MM-DD` representing a date of scheduled MLB games.


## Local Development
### Formatting
Format the source code

```make fmt```
### Test
Run unit tests

```make test```

### Run
Run the backend service via `make` or `docker`

```make run```

or

```
docker build -t mlbtakehome .
docker run -it --rm -p 8080:8080 mlbtakehome
```

Running either of these commands will serve the `/schedule` API at the following endpoint: `localhost:8080/api/v1/schedule?teamId=<teamId>&date=<YYYY-MM-DD>`.  

Example:
```
curl 'localhost:8080/api/v1/schedule?teamId=141&date=2021-09-11'
```


