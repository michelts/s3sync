# S3Sync

This project allows migrating files between Oracle and AWS object storage
clouds, in order to keep compatibility of a legacy app in the process of being
replaced.

To run the project, build the docker image and execute the command below:

```
docker-compose build app
docker-compose run app
go run main.go
```

Before that, you might provide a `.env` file similar to the one below:

```
BUCKETS=bucket-a,bucket-b

# Origin
SOURCE_ACCESS_KEY=AWS-access-key
SOURCE_SECRET_KEY=AWS-secret-key
SOURCE_REGION=AWS-region-name

# Destination
DEST_ACCESS_KEY=ORACLE-access-key
DEST_SECRET_KEY=ORACLE-secret-key
DEST_ENDPOINT_URL=ORACLE-endpoint-url
DEST_REGION=ORACLE-region-name
```
