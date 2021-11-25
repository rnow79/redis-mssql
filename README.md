# redis-mssql
Redis power demonstration! This is a simple Go API for demonstrate how redis can help serving images. It can retrieve one image, or it can retrieve many (one by one, simulating calls on the same web page). Only one connection to SQL is made, and only one to Redis. The image, by default, have a 10px red border when retrieved from SQL and a 10px green border when retrieved from Redis. Also there is a default 10 sec TTL for all Redis registers. The API shows how many time is spent in retrieving the images (if you get 10 images from SQL and then refresh the webpage before 10 secs, you will see the time diff). Redis Rocks!

## Set MSSQL variables
First, you have to put your own variables:

* <code>var mssqlServer = "myserver"</code> IP or host of the MSSQL server
* <code>var mssqlPort = 1433</code> MSSQL port
* <code>var mssqlUser = "myusername"</code> MSSQL username
* <code>var mssqlPassword = "mypassword"</code> Gues what... MSSQL user foot size
* <code>var mssqlSelect = "select photo from mydatabase.dbo.photos where id='%s'"</code> Ensure that your SQL select returns only one row, and is of type image.

## Set REDIS variables
* <code>Addr</code> Redis host.
* <code>Password</code> Redis password.
* <code>DB</code> Redis Databalse.
* <code>var keyExpiration time.Duration = 10 * time.Second</code> Expiration time on Redis database for the images

## Run the project
You can either run the project locally, or create a Docker image (Dockerfile included)

## Open navigator
<code>http://localhost/id=5000</code> Retrieves the image with id 5000

<code>http://localhost/id=5000|5001|5002</code> Simulates three calls to the database. If any of the ids is found in Redis, the API serves it. Otherwise, it's served from the SQL Server. You can concatenate the ids you want. Simple quotes are removed from ids for preventing SQL injection, and if a particular id is not found in the database, the API will ignore that picture.

## Use
This API is only for demonstration purpose.



