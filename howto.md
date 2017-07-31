# How To Interact With Logisense Internal Server #

- Logisense performance server is hosted by qa-performance VM and is listening on port 4999.
- All steps below are to give a brief introduction, for detailed description for each endpoint, refer to [Endpoints](https://github.com/han-hgu/perf-prototype/blob/master/readme.md#endpoints) section.
- Testrun information is searchable through tags. For example, run the following request to retrieve all tests with tag "8.5.26.5" and "Rating", only the testID and metadata of the test is returned to give you an idea of what the testrun is about. 

	**Request**

		GET http://qa-performance:4999/v1/tests?tag=8.5.26.5&tag=Rating

	**Returned two tests, tests returned have both tags**

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

- Run the following request if you want to see the completed stats collected for the test run. The below request returns you the full result for testID "596901ee8a8fd50ac009ab7e". Result is not pasted here as it is a much bigger dataset with all aspect of the application server and database KPIs collected in real time. 

	**Request**

		GET http://qa-performance:4999/v1/tests/596901ee8a8fd50ac009ab7e
	
- The framework also provides a way to visualize the KPI comparison among different rating test runs. For example, the following request compares the three testruns in
	- The Momentum 8.5.0.13 production system (id: 596925d48a8fd50b58ff09fd)
	- The upgraded system with EIP 8.5.26.5-Hotfix.6RC24 installed (id: 5965287b8a8fd50bb432e834)
	- The upgraded system tuned for better performance(TS-4240) with EIP 8.5.26.5-Hotfix.6RC28 installed (id: 596901ee8a8fd50ac009ab7e)

	All three runs rate against a fixed Momentum dataset.


	**Request**

		GET http://qa-performance:4999/v1/rating/charts?id=596925d48a8fd50b58ff09fd&id=5965287b8a8fd50bb432e834&id=596901ee8a8fd50ac009ab7e


	**Returns**
	![IBBS Rating Perf](https://user-images.githubusercontent.com/12279676/28645017-db739f28-7229-11e7-9c9d-44bcf5fd028d.png)
