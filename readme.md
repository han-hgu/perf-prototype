# Perf
EngageIP performance testing framework.

## Getting Started
Under development.

## Building
You need a Go development to build the executable. To install all dependencies, change directory into the project and run:

    go get ./...

You can cross compile it using 
`GOOS` and `GOARCH`

    GOOS=windows GOARCH=amd64 go build

## Running
- Running the executable starts up a server listening on port 4999

## Endpoints

###Rating
- **[<code>GET</code> tests](https://github.com/han-hgu/perf-prototype/blob/master/api-documentation/rating/GET_tests.md)**
- **[<code>GET</code> tests/:id](https://github.com/han-hgu/perf-prototype/blob/master/api-documentation/rating/GET_tests_id.md)**
- **[<code>POST</code> tests](https://github.com/han-hgu/perf-prototype/blob/master/api-documentation/rating/POST_tests.md)**

###Billing


## FAQ ##


## Requirements ##

- Rating perf testing
- Billing perf testing
- File drop service to generate and create the input files in the designed location
- Controllable rate for file drop
- Rates are captured through out the whole process not just an average value
- Error handling:
	- Server recovery
	- Error stats
	- Persistent storage

[OAuth]: http://oauth.net/core/1.0a/
[Beginner’s Guide]: http://hueniverse.com/oauth/
[JSON]: http://json.org
[quick tutorial]: http://www.webmonkey.com/2010/02/get_started_with_json/
[Register your application]: http://500px.com/settings/applications
[API Terms of Use]: https://github.com/500px/api-documentation/blob/master/basics/terms_of_use.md
[See if the concepts used by the API are familiar to you]: https://github.com/500px/api-documentation#what-do-i-need-to-know-before-i-start-using-the-api
[A good md reference page for api]:
https://github.com/500px/api-documentation/blob/master/README.md
