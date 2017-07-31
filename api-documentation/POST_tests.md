# Start a test

    POST tests

## Description
Start monitoring a rating/billing test

***

## Parameters

- **type** - Test type, either "rating" or "billing"

***

## Request payload
Payload in JSON format with the following key/value pairs

- **db_config** - [Database Server Configuration](), JSON dictionary
- **app_config** - [Application Server Configuration](), JSON dictionary
- **chart_config** - [Chart Configuration](), JSON dictionary
- **comment** - Additional comment for the testrun, string
- **tags** - A list of tags associated with the testrun, JSON array
- **collection_interval** - How often the framework tracks KPIs, default is "30s" if not specified; A valid input is a sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "1.5h" or "2h45m", valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"

### Additional payload for rating

- **use\_existing\_file** - If true the system will not generate rating input files, false otherwise
- **number\_of\_files** - The number of rating input files waiting to be processed for the testrun. The framework will consult the rating system for how many input files has been processed so far; If the number of files processed is equal to this number, the framework marks the test as completed. If 
**use\_existing\_file** is set to false, the framework is responsible for generating this amount of input files according to the parameters below
- **raw\_fields** - A raw record as a reference record for the framework to generate the input files, JSON array, only used if **use\_existing\_file** is set to false
- **amount\_field\_index** - The index of the amount field in **raw\_fields** array, the framework uses this to insert different UDR amount values, only used if **use\_existing_file** is set to false
- **timestamp\_field\_index** - The index of the timestamp field in **raw\_fields** array, only used if **use\_existing\_file** is set to false
- **number\_of\_records\_per\_file** - The number of raw records generated in each input file, only used if **use\_existing\_file** is set to false
- **drop\_location** - The location where the input files are created by the framework, only used if **use\_existing\_file** is set to false
- **filename\_prefix** - The prefix used for the name of input files, default to TestID if not specified, only used if **use\_existing\_file** is set to false

### Additional payload for billing

- **owner_name** - The name of the owner under which bill run is started


***

## Return format
A map with the following keys and values:

- **id** - Test ID


***

## Errors

- **400 Bad Request** — Invalid/missing test type


***

## Example
**Request**

    POST v1/test?type=rating

**Body**

	{
		"db_config": {
			"ip": "192.168.1.47",
			"port": 1433,
			"db_name": "EngageIP_NonRevenue",
			"uid": "sa",
			"password": "Q@te$t#1",
			"perfmon_url": "http://192.168.1.47:5000/v1"
		},
	
		"app_config": {
			"version": "EngageIP 8.5.26.5-Hotfix.6RC24",
			"perfmon_url": "http://192.168.1.51:5000/v1"
		},
	
		"chart_config": {
			"title": "8.5.26.5-Hotfix.6RC24"
		},
	
		"tags": [
			"Momentum",
			"rating",
			"A",
			"B"
		],
	
		"use_existing_file": true,
		"collection_interval": "300s",
		"number_of_files": 3000
	}


**Return** __shortened for example purpose__

	{
		"id": "597f73de8a8fd509a801e16a"
	}
	


[photo stream]: https://github.com/500px/api-documentation/blob/master/basics/formats_and_terms.md#500px-photo-terms
[OAuth]: https://github.com/500px/api-documentation/tree/master/authentication
[http://500px.com/:username]: http://500px.com/iansobolev
[http://500px.com/:username/following]: http://500px.com/iansobolev/following
[category]: https://github.com/500px/api-documentation/blob/master/basics/formats_and_terms.md#categories
[short format]: https://github.com/500px/api-documentation/blob/master/basics/formats_and_terms.md#short-format-1
[photo sizes]: https://github.com/500px/api-documentation/blob/master/basics/formats_and_terms.md#image-urls-and-image-sizes
