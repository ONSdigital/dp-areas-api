dp-areas-api
================
Digital Publishing areas API used to navigate profiles for different geographical areas across the UK

### Getting started

* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default   | Description
| ---------------------------- | --------- | -----------
| BIND_ADDR                    | :25500    | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s        | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL         | 30s       | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s       | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)

### Connecting to the AWS AURORA RDS instance from your local machine

Note: `<RDS_INSTANCE_ENDPOINT>` can be obtained from the AWS AURORA RDS cluster console here:

https://eu-west-1.console.aws.amazon.com/rds/home?region=eu-west-1#database:id=develop-area-profiles-postgres;is-cluster=true

1. Add DB host to /etc/hosts:

`
127.0.0.1 <RDS_INSTANCE_ENDPOINT>
`

2. Set the following in your environment - remote postgres connection:

```
export AWS_PROFILE="development"
export PG_USER="dp-areas-api-publishing"
export PGPASSWORD="$(aws rds generate-db-auth-token --hostname <RDS_INSTANCE_ENDPOINT> --port 5432 --region <AWS_REGION> --username dp-areas-api-publishing)"
```

These config variables are only required if running the dp-areas-api branch I’m working on (https://github.com/ONSdigital/dp-areas-api/tree/feature/postgres_healthcheck):

```
export DBNAME="dp-areas-api"
export DBUSER="dp-areas-api-publishing"
export DBHOST="<RDS_INSTANCE_ENDPOINT>"
export DBPORT=5432
export AWSREGION=<AWS_REGION>
```

for local postgres connection (relies on `dp-compose`):

*Note:* set _*DPPostgresLocal*_ to _*true*_ to use *local postgres instance*

```
export DPPostgresLocal=true
export DPPostgresUserName="postgres"`
export DPPostgresUserPassword="<PASSWORD>"` (see docker compose)
export DPPostgresLocalPort="5432"`
export DPPostgresLocalDB="dp-areas-api"`
```

3. Run the dp command:

```
dp remote allow develop
```

4. Open a port forwarding connection from your localhost:5432 to the AWS AURORA RDS instance endpoint by running:

```
dp ssh develop publishing 3 -v -- -L 5432:<RDS_INSTANCE_ENDPOINT>:5432
```

5. Get the required certificate by running:

```
wget https://s3.amazonaws.com/rds-downloads/rds-ca-2019-root.pem
```

6. Finally, execute the sql command to open a tunnel to the AWS AURORA RDS instance:

```
psql -h <RDS_INSTANCE_ENDPOINT> -p 5432 "sslmode=verify-full sslrootcert=<PATH_TO_PEM_FILE> dbname=dp-areas-api user=dp-areas-api-publishing"
```

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2021, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

