#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" <<-EOSQL
    -- Create mattermost database if it doesn't exist
    SELECT 'CREATE DATABASE mattermost'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'mattermost')\gexec
    
    -- Create keycloak database if it doesn't exist
    SELECT 'CREATE DATABASE keycloak'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'keycloak')\gexec
    
    -- Grant privileges on mattermost database
    GRANT ALL PRIVILEGES ON DATABASE mattermost TO "$POSTGRES_USER";
    
    -- Grant privileges on keycloak database
    GRANT ALL PRIVILEGES ON DATABASE keycloak TO "$POSTGRES_USER";
EOSQL

echo "Databases initialized successfully"
