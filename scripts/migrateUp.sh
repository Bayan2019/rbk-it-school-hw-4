#!/bin/bash

if [ -f .env ]; then
    source .env
fi


cd internal/repository/postgres/sql/schema
goose postgres $DATABASE_URL up