# Drawing Charts for KPIs

	GET rating/charts



## Description
Draw charts for all KPIs captured. KPIs for multiple test runs could be drawn in the same chart for comparison purpose

***
## Parameters
- **id** - Can be specified multiple times for KPI comparison for different test runs; test runs must have the same type and `collection interval` value

***

## Chart returned
- **UDR Rates** - UDR rate per interval
- **Total UDR processed** - Total UDR records generated per interval
- **Application Server CPU utilization** - CPU utilization for the application server per interval 
- **Application Server Memory utilization** - Memory utilization for the application server per interval 
- **Database Server CPU utilization** - CPU utilization for the database `db_name` per interval 
- **Database Logical Reads** - Logical reads for database `db_name` per interval
- **Database Physical Reads** - Physical reads for database `db_name` per interval
- **Database Logical Writes** - Logical writes for database `db_name` per interval
