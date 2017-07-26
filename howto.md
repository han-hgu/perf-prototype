# How To Interact With Logisense Internal Server #

- Internal performance server is hosted by qa-performance VM and listening on port 4999.
- All the following steps are to give a brief introduction, for detailed description of each endpoint, refer to [Endpoint](readme#Endpoints) section.
- The meta To retrieve all test meta-data containing "8.5.26.5" and "Rating" tags
		
	**Request**
			
		GET http://qa-performance:4999/v1/tests?tag=8.5.26.5&tag=Rating

	**Returned two tests, only tests with both tags are returned**
			
		[
			{
				"id": "5965287b8a8fd50bb432e834",
				"meta_data": 
				{
					"test_type": "rating",
					"start_date": "2017-07-11T15:35:23.577-04:00",
					"test_completed": true,
					"test_duration": "6h34m59.3259932s",
					"tags": [
						"8.5.26.5",
					 	"Momentum",
					 	"Rating",
					 	"TS-4240",
						"TS-3038",
						"Hotfix.6",
						"Broadworks"
					],
					"collection_interval": "30s",
					"app_param": {...},
					"db_param": {...},
					"chart_title": "8.5.26.5-Hotfix.6RC24"
				}
			},
			{
				"id": "596901ee8a8fd50ac009ab7e",
				"meta_data": {
					"test_type": "rating",
					"start_date": "2017-07-14T13:39:58.613-04:00",
					"test_completed": true,
					"test_duration": "32m15.9656547s",
					"tags": [
					"8.5.26.5",
					"Momentum",
					"Rating",
					"Hotfix.6",
					"Broadworks",
					"TS-4240",
					"TS-3038",
					"TS-5139"
				],
				"comment": "Run 'exec [Support_AddBillersBucketsBulk] '20161101''",
				"collection_interval": "30s",
				"app_param": {},
				"db_param": {},
				"chart_title": "8.5.26.5-Hotfix.6RC28"
				}
			}
		]

- Run the following request to retrieve the full result using the test ID; The following request returns you the full test result with id "596901ee8a8fd50ac009ab7e", it is a much bigger dataset with all stats collected during the test run

	**Request**
	
	```
	GET http://qa-performance:4999/v1/tests/596901ee8a8fd50ac009ab7e
	```

- The framework provides a way to visually compare system KPIs among different rating test runs. For example, the following query compares the test run simulating the Momentum 8.5.0.13 production system(id: 596925d48a8fd50b58ff09fd), the upgraded system with stock EIP 8.5.26.5-Hotfix.6RC24 installed(id: 5965287b8a8fd50bb432e834), and the one tuned for better performance(TS-4240) with version 8.5.26.5-Hotfix.6RC28 installed(id: 596901ee8a8fd50ac009ab7e). All three runs rate against a fixed dataset tailored for Momentum. 

	**Request**

	```
	GET http://qa-performance:4999/v1/rating/charts?id=596925d48a8fd50b58ff09fd&id=5965287b8a8fd50bb432e834&id=596901ee8a8fd50ac009ab7e
	```

	**Returns**
	![IBBS Rating Perf](https://github.com/han-hgu/perf-prototype/blob/master/api-documentation/rating/IBBSRatingPerfComparison.png)

	
	