#!/bin/bash
set -e

POSTGRES_USER="${POSTGRES_USER:-feedback}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-feedback}"
PGHOST="${PGHOST:-db}"
PGPORT="${PGPORT:-5432}"

export PGPASSWORD="${POSTGRES_PASSWORD}"

echo "Waiting for PostgreSQL to be ready..."
until pg_isready -h "${PGHOST}" -p "${PGPORT}" -U "${POSTGRES_USER}" -d postgres; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done

echo "PostgreSQL is ready - creating databases..."

# Create databases if they don't exist
psql -v ON_ERROR_STOP=1 -h "${PGHOST}" -p "${PGPORT}" -U "${POSTGRES_USER}" -d postgres <<-EOSQL
    -- Create mattermost database if it doesn't exist
    SELECT 'CREATE DATABASE mattermost'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'mattermost')\gexec
    
    -- Create keycloak database if it doesn't exist
    SELECT 'CREATE DATABASE keycloak'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'keycloak')\gexec
    
    -- Grant privileges on mattermost database
    GRANT ALL PRIVILEGES ON DATABASE mattermost TO "${POSTGRES_USER}";
    
    -- Grant privileges on keycloak database
    GRANT ALL PRIVILEGES ON DATABASE keycloak TO "${POSTGRES_USER}";
EOSQL

echo "Databases initialized successfully"

