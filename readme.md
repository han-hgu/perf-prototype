# EIP Performance Test Framework
The framework is responsible for monitoring test progress, collecting real-time KPIs and managing test results.
The API is REST API, return format for all endpoints is JSON.

[How To Interact With Logisense Internal Server](https://github.com/han-hgu/perf-prototype/blob/master/howto.md)

## Endpoints
### Querying Test Results
- **[<code>GET</code> tests](https://github.com/han-hgu/perf-prototype/blob/master/api-documentation/GET_tests.md)**
- **[<code>GET</code> tests/:id](https://github.com/han-hgu/perf-prototype/blob/master/api-documentation/GET_tests_id.md)**

### Posting Tests
- **[<code>POST</code> tests](https://github.com/han-hgu/perf-prototype/blob/master/api-documentation/rating/POST_tests.md)**

### Drawing Comparison Charts
- **[<code>GET</code> rating/charts](https://github.com/han-hgu/perf-prototype/blob/master/api-documentation/rating/GET_rating_charts.md)**


## Building
You need a Go development to build the executable. To install all dependencies, change directory into the project and run:

    go get ./...

You can cross compile it using
`GOOS` and `GOARCH`, to build a Windows 64-bit version for example:

    GOOS=windows GOARCH=amd64 go build

## Running
- Running the executable starts up a server on port 4999

## Requirements
-

## Design Decisions
### Database
- [MongoDB](https://www.mongodb.com/) for saving test results
- [MongoDB Go driver](https://labix.org/mgo)

### Cache
- [BigCache](https://github.com/allegro/bigcache)

### Graphing
- [Google Charts](https://developers.google.com/chart/)

## TODO
- Use [Context package](https://golang.org/pkg/context/)
- Authentication
- Controllable rate for file drop
- Error handling:
	- Service recovery
	- Error stats

[OAuth]: http://oauth.net/core/1.0a/
[Beginner’s Guide]: http://hueniverse.com/oauth/
[JSON]: http://json.org
[quick tutorial]: http://www.webmonkey.com/2010/02/get_started_with_json/
[A good md reference page for api]: https://github.com/500px/api-documentation/blob/master/README.md
