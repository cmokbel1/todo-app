#!/bin/bash

set -e

echo "POSTGRES_USER=${POSTGRES_USER}"
echo "DB_USER=${DB_USER}"
echo "DB_NAME=${DB_NAME}"

function create_database_and_user() {
	local db=$1
	# Create the postgres if it doesn't exist
	psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" <<-EOSQL
		SELECT 'CREATE DATABASE ${db}'
		WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '${db}')\gexec
EOSQL

	psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" --dbname="${db}" <<-EOSQL
	    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    	CREATE EXTENSION IF NOT EXISTS "pgcrypto";
    	CREATE EXTENSION IF NOT EXISTS "hstore";
EOSQL
}

psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" <<-EOSQL
	    CREATE USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}';
		ALTER USER ${DB_USER} with SUPERUSER;
EOSQL

if [[ -n "${DB_NAME}" ]]; then
	for db in $(echo ${DB_NAME} | tr ',' ' '); do
		echo "Creating database: '${db}'"
		create_database_and_user ${db}
	done
else
	echo "Missing DB_NAME env variable"
	exit 1
fi