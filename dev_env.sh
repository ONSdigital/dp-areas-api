#!/bin/bash
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=
export DBNAME="dp-areas-api"
export DBUSER="dp-areas-api-publishing"
export DBHOST=<RDS_INSTANCE_ENDPOINT>
export DBPORT=5432
export AWSREGION=eu-west-1
export PG_USER="dp-areas-api-publishing"
export PGPASSWORD="$(aws rds generate-db-auth-token --hostname <RDS_INSTANCE_ENDPOINT> --port 5432 --region eu-west-1 --username dp-areas-api-publishing)"