# Query Test Metadata

	GET tests

## Description
Returns the metadata portion of all the tests(in-progress and completed), which contains every tag specified in the request. If no tag is specified, all test metadata are returned  

***
## Parameters
- **tag** - Can be specified multiple times, only test metadata containing every tag specified are returned, case sensitive

## Return format
An array of objects with the following keys and values:

- **id** - Test ID. It can be used to query the complete test result
- **meta_data** - [Metadata](https://github.com/han-hgu/perf-prototype/blob/master/api-documentation/terms.md#metadata) describing the test run

***

## Example

**Request**

	GET v1/tests?tag=A&tag=B

## Return

	[{
		"id": "597a564b8a8fd509a801e169",
		"meta_data": {
			"test_type": "rating",
			"start_date": "2017-07-27T17:08:27.7991098-04:00",
			"test_completed": false,
			"tags": [
				"Momentum",
				"rating",
				"A",
				"B"
			],
			"collection_interval": "300s",
			"app_param": {
				"version": "EngageIP 8.5.26.5-Hotfix.6RC24",
				"EIP_option": null,
				"perfmon_url": "http://192.168.1.51:5000/v1",
				"sys_info": {
					"cpu_info": [{
						"cacheSize": 0,
						"coreId": "",
						"cores": 8,
						"cpu": 0,
						"family": "179",
						"flags": [],
						"mhz": 2500,
						"microcode": "",
						"model": "",
						"modelName": "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
						"physicalId": "0000000000000000",
						"stepping": 0,
						"vendorId": "GenuineIntel"
					}],
					"host_info": {
						"bootTime": 1500284637,
						"hostid": "39f2bf70-0eb6-4362-be9c-c497cb19a933",
						"hostname": "QA-Mommentum-App-Server",
						"kernelVersion": "",
						"os": "windows",
						"platform": "Microsoft Windows Server 2012 R2 Datacenter",
						"platformFamily": "Server",
						"platformVersion": "6.3.9600 Build 9600",
						"procs": 55,
						"uptime": 905071,
						"virtualizationRole": "",
						"virtualizationSystem": ""
					},
					"mem_info": {
						"active": 0,
						"available": 91146166272,
						"buffers": 0,
						"cached": 0,
						"dirty": 0,
						"free": 0,
						"inactive": 0,
						"pagetables": 0,
						"shared": 0,
						"slab": 0,
						"swapcached": 0,
						"total": 94368731136,
						"used": 3222564864,
						"usedPercent": 3,
						"wired": 0,
						"writeback": 0,
						"writebacktmp": 0
					}
				}
			},
			"db_param": {
				"db_name": "EngageIP_NonRevenue",
				"perfmon_url": "http://192.168.1.47:5000/v1",
				"compatibility_level": 110,
				"sys_info": {
					"cpu_info": [{
							"cacheSize": 0,
							"coreId": "",
							"cores": 24,
							"cpu": 0,
							"family": "179",
							"flags": [],
							"mhz": 2500,
							"microcode": "",
							"model": "",
							"modelName": "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
							"physicalId": "BFEBFBFF000306F2",
							"stepping": 0,
							"vendorId": "GenuineIntel"
						},
						{
							"cacheSize": 0,
							"coreId": "",
							"cores": 24,
							"cpu": 1,
							"family": "179",
							"flags": [],
							"mhz": 2500,
							"microcode": "",
							"model": "",
							"modelName": "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
							"physicalId": "BFEBFBFF000306F2",
							"stepping": 0,
							"vendorId": "GenuineIntel"
						}
					],
					"host_info": {
						"bootTime": 1500432852,
						"hostid": "0b2d762b-c791-419b-a752-2e1f9bc0cbc3",
						"hostname": "BMSQL-HOST-QA01",
						"kernelVersion": "",
						"os": "windows",
						"platform": "Microsoft Windows Server 2012 R2 Standard",
						"platformFamily": "Server",
						"platformVersion": "6.3.9600 Build 9600",
						"procs": 76,
						"uptime": 756863,
						"virtualizationRole": "",
						"virtualizationSystem": ""
					},
					"mem_info": {
						"active": 0,
						"available": 17765199872,
						"buffers": 0,
						"cached": 0,
						"dirty": 0,
						"free": 0,
						"inactive": 0,
						"pagetables": 0,
						"shared": 0,
						"slab": 0,
						"swapcached": 0,
						"total": 274777624576,
						"used": 257012424704,
						"usedPercent": 93,
						"wired": 0,
						"writeback": 0,
						"writebacktmp": 0
					}
				}
			},
			"chart_title": "8.5.26.5-Hotfix.6RC24"
		}
	}]
