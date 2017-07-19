# Rating

    POST /rating/tests

## Description
Create a rating test
***

## Requires authentication

***

## Parameters

***

## Request body format
A map with the following keys and values:

- **amount\_field\_index** - The amount field index in **raw_fields**
- **timestamp\_field\_index** - The timestamp field index in **raw_fields**
- **use\_existing_file** - If true the system will not generate rating input files
- **number\_of\_files** - Number of rating input files to generate if **use\_existing_file** is set to false
- **number\_of_records\_per_file** - Number of raw records to generate if **use\_existing_file** is set to false
- **raw\_fields** - A raw record for reference while generating input files
- **drop\_location** - The location where the input files are generated
- **filename\_prefix** - The filename prefix used for input files, default to use UUID if not specified
- **starting\_id** - The system ignores the event log if its id is less than this. If not specified, **starting\_id** is default to 0
- **additional\_info** - A map that will be appended to the test result info




***

## Return format
An array with the following keys and values:

- **feature** — Feature that is being returned.
- **filters** — Additional filters that were used:
    - 'category' — The ID of the **[category][]** that photos were filtered by;
    - 'user_id' — The ID of the user specified by 'user_id' or 'username' parameters;
    - 'friends_ids' — IDs of users the user specified is following;
- **current_page** — Number of the page that is returned.
- **total_pages** — Total number of pages in this feature's stream.
- **total_items** — Total number of items in this feature's stream.
- **photos** — An array of Photo objects in **[short format](https://github.com/500px/api-documentation/blob/master/basics/formats_and_terms.md#short-format)**.

***

## Errors

***

## Example
**Request**

    POST v1/test?type=billing

**Body**
```json
{  
	"db_config":{
		"ip": "192.168.1.47",
		"port": 1433,
		"db_name": "EngageIP_Revenue",
		"uid": "sa",
		"password": "Q@te$t#1",
		"perfmon_url": "http://192.168.1.47:5000/v1/stats",
		"additional_info": {
			"CPU": "2 x Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz, 2500 Mhz, 12 Core(s), 24 Logical Processor(s)",
			"RAM": "256 GB",
			"OS": "Microsoft Windows Server 2012 R2 Standard",
			"system_model": "PowerEdge R730xd",
			"comment": "Specs from Windows System Information"
		}
	},

	"app_config": {
		"version": "EngageIP 8.5.26.5-Hotfix.6RC28",
		"perfmon_url": "http://192.168.1.51:5000/v1/stats",
		"EIP_option": [],
		"additional_info": {
			"CPU": "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz, 2500 Mhz, 8 Core(s), 8 Logical Processor(s)",
			"RAM": "85.9 GB",
			"VM": "True",
			"OS": "Microsoft Windows Server 2012 R2 Datacenter",
			"system_model": "Virtual Machine",
			"comment": "Specs from Windows System Information"
		}
	},

	"comment": "Ran updated UDRBillerFixer from TS-4942",

	"chart_config": {
		"title": "8.5.26.5-Hotfix.6RC28"
	},

	"tags": [
		"8.5.26.5",
		"Momentum",
		"billing",
		"Hotfix.6",
		"TS-4942",
		"Momentum_Voice"
	],

    "owner_name": "Momentum_Voice",
    "collection_interval": "10s"
}
```

**Return** __shortened for example purpose__
``` json
{
  "feature": "popular",
  "filters": {
      "category": false,
      "exclude": false
  },
  "current_page": 1,
  "total_pages": 250,
  "total_items": 5000,
  "photos": [
    {
      "id": 4910421,
      "name": "Orange or lemon",
      "description": "",
      "times_viewed": 709,
      "rating": 97.4,
      "created_at": "2012-02-09T02:27:16-05:00",
      "category": 0,
      "privacy": false,
      "width": 472,
      "height": 709,
      "votes_count": 88,
      "comments_count": 58,
      "nsfw": false,
      "image_url": "http://pcdn.500px.net/4910421/c4a10b46e857e33ed2df35749858a7e45690dae7/2.jpg",
      "user": {
        "id": 386047,
        "username": "Lluisdeharo",
        "firstname": "Lluis ",
        "lastname": "de Haro Sanchez",
        "city": "Sabadell",
        "country": "Catalunya",
        "fullname": "Lluis de Haro Sanchez",
        "userpic_url": "http://acdn.500px.net/386047/f76ed05530afec6d1d0bd985b98a91ce0ce49049/1.jpg?0",
        "upgrade_status": 0
      }
    },
    {
      "id": 4905955,
      "name": "R E S I G N E D",
      "description": "From the past of Tagus River, we have History and memories, some of them abandoned and disclaimed in their margins ...",
      "times_viewed": 842,
      "rating": 97.4,
      "created_at": "2012-02-08T19:00:13-05:00",
      "category": 0,
      "privacy": false,
      "width": 750,
      "height": 500,
      "votes_count": 69,
      "comments_count": 29,
      "nsfw": false,
      "image_url": "http://pcdn.500px.net/4905955/7e1a6be3d8319b3b7357c6390289b20c16a26111/2.jpg",
      "user": {
        "id": 350662,
        "username": "cresendephotography",
        "firstname": "Carlos",
        "lastname": "Resende",
        "city": "Forte da Casa",
        "country": "Portugal",
        "fullname": "Carlos Resende",
        "userpic_url": "http://acdn.500px.net/350662.jpg",
        "upgrade_status": 0
      }
    }
  ]
}
```

[photo stream]: https://github.com/500px/api-documentation/blob/master/basics/formats_and_terms.md#500px-photo-terms
[OAuth]: https://github.com/500px/api-documentation/tree/master/authentication
[http://500px.com/:username]: http://500px.com/iansobolev
[http://500px.com/:username/following]: http://500px.com/iansobolev/following
[category]: https://github.com/500px/api-documentation/blob/master/basics/formats_and_terms.md#categories
[short format]: https://github.com/500px/api-documentation/blob/master/basics/formats_and_terms.md#short-format-1
[photo sizes]: https://github.com/500px/api-documentation/blob/master/basics/formats_and_terms.md#image-urls-and-image-sizes
