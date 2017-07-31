# Terms
## Metadata
Only in response, the metadata of a testrun includes the following

- **type** - Type of the test, either "rating" or "billing"
- **start_date** - Start time of the test
- **test_completed** - Whether the test is completed, boolean
- **duration** - Duration of the test if completed
- **tags** - A list of tags associated with the test, copied from request
- **comment** - Additional comment for the test, string
- **collection_interval** - The framework collects the stats every `collection interval`
- **app_param** - The software/hardware information of the application server, see [Application Server Parameters](#application-server-parameters) for details
- **db_param** - The software/hardware information of the database server, see [Database Server Parameters](#database-server-parameters) for details
- **chart_title** - The title of the testrun in comparison charts




## Application Server Configuration
Only in request

- **version** - EngageIP version installed
- **EIP_option** - A list of EngageIP options for the testrun, JSON list
- **perfmon_url** - URL of the performance monitor for KPIs and specs
- **additional_info** - Additional information about the application server for the user to keep track of, this needs to be entered manually and will be saved in [Application Server Parameters](#application-server-parameters), JSON dictionary



## Application Server Parameters
Only in response

- **version** - EngageIP version installed, this is copied from the same field in [Application Server Configuration](#application-server-configuration)
- **EIP_option** - A list of EngageIP options for the testrun, this is copied from the same field in [Application Server Configuration](#application-server-configuration)
- **perfmon_url** - URL of the performance monitor for KPIs and specs, this is copied from the same field in [Application Server Configuration](#application-server-configuration)
- **sys_info** - Hardware information of the application server, auto-detected, JSON dictionary
- **additional_info** - Additional information about the application server, this is copied from the same field in [Application Server Configuration](#application-server-configuration)


## Application Server Statistics
Only in response

- **cpu(%)** - CPU usage is collected every `collection_interval` throughout the test, stored as a list of float values
- **cpu_max(%)** - The peak CPU usage throughout the test
- **mem(%)** - Memory usage in percentage collected every `collection_interval`, stored as a list of float values
- **mem_max(%)** - The peak memory usage in percentage throughout the test


## Database Server Configuration
Only in request

- **ip** - IP address of the database server
- **port** - Port of the database server
- **db_name** - Name of the database
- **uid** - UID to connect to the database server, for example "sa"
- **password** - Password for UID
- **perfmon_url** - URL of the performance monitor for KPIs and specs
- **additional_info** - Additional information about the database server, this needs to be manually entered in request in JSON dictionary format and will be copied to [Database Server Parameters](#database-server-parameters)


## Database Server Parameters
Only in response

- **db_name** - The name of the database, this is copied from the same field in [Database Server Configuration](#database-server-configuration)
- **perfmon_url** - URL of the performance monitor for KPIs and specs, this is copied from the same field in [Database Server Configuration](#database-server-configuration)
- **compatibility_level** - The compatibility mode of the database, auto-detected
- **sys_info** - Hardware information of the database server, auto-detected, JSON dictionary
- **additional_info** - Additional information about the database server, this is copied from the same field in [Database Server Configuration](#database-server-configuration)

## Database Server Statistics
Only in response

- **cpu(%)** - CPU usage is NOT the DB server overall usage but the usage of the specific database according to `db_name`, collected every `collection_interval`, stored as a list of float values
- **cpu_max(%)** - Peak CPU usage of the specific database
- **mem(%)** - Memory usage in percentage of the specific database, collected every `collection_interval`, stored as a list of float values
- **mem_max(%)** - The peak memory usage of the specific database throughout the test
- **logical\_reads\_total** - Total logical reads performed in the specific  database throughout the test
- **logical_reads** - An array of logical read samples collected every `collection_interval` for the specific database throughout the test
- **logical\_writes\_total** - Total logical writes performed in the specific  database throughout the test
- **logical_writes** - An array of logical write samples collected every `collection_interval` for the specific database throughout the test
- **physical\_reads\_total** - Total physical reads performed in the specific  database throughout the test
- **physical_reads** - An array of physical read samples collected every `collection_interval` for the specific database throughout the test



## Chart Configuration


- **title** - The title of the testrun in KPI comparison charts
