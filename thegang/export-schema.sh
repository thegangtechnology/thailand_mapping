#!/bin/bash

# Set the container name
container_name="go-backend-postgres-1"
user="go-backend-postgres-1"
database="go-backend-postgres-1"

# Export table schemas to a csv file
echo "Exporting table schemas to a csv file..."
docker exec $container_name psql -U $user -d $database -c "\copy (SELECT table_name, column_name, data_type, is_nullable FROM information_schema.columns WHERE table_schema='public' order by table_name) TO '/tmp/table_schemas.csv' with csv"

# Copy the csv file from the container to local machine
echo "Copying the csv file from the container to local machine..."
docker cp $container_name:/tmp/table_schemas.csv ./thegang/manual