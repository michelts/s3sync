# S3Sync

This project allows migrating files between Oracle and AWS object storage
clouds, in order to keep compatibility of a legacy app in the process of being
replaced.

To run the project, you can use the dev image. For production usage, build the
prod image and start it. The service would be available in the url
http://localhost:5000/ as a POST request expecting a json payload containing a
`Publisher`, `Publication` and `Issue` ids, all integers.

Before starting the service, you might provide a `.env` file similar to the one
below:

```
OCI_ACCESS_KEY=ORACLE-access-key
OCI_SECRET_KEY=ORACLE-secret-key
OCI_ENDPOINT_URL=ORACLE-endpoint-url
OCI_REGION=ORACLE-region-name

AWS_ACCESS_KEY=AWS-access-key
AWS_SECRET_KEY=AWS-secret-key
AWS_REGION=AWS-region-name
```

To run the application manually, there's a helper script `batch_exec.py` that
expects a csv file `items.csv` in the same directory. The csv file must contain
the `Publisher`, `Publication` and `Issue` ids, on item per row.
