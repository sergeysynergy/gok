#!/bin/bash

# Below this comment script will insert current user name while building: user=$USER
# Below this comment script will insert user password while building: password=$USER_PASSWORD

psql -U postgres <<- EOSQL
    CREATE USER $user WITH PASSWORD '$password';

    CREATE DATABASE "storage";
    GRANT ALL PRIVILEGES ON DATABASE "storage" TO $user;

    CREATE DATABASE "storage_test";
    GRANT ALL PRIVILEGES ON DATABASE "storage_test" TO $user;

    CREATE DATABASE "auth";
    GRANT ALL PRIVILEGES ON DATABASE "auth" TO $user;

    CREATE DATABASE "auth_test";
    GRANT ALL PRIVILEGES ON DATABASE "auth_test" TO $user;
EOSQL
