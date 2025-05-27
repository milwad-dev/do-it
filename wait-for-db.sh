#!/bin/sh

echo "Waiting for MySQL to be ready at $DB_HOST:$DB_PORT..."

while ! mysqladmin ping -h"$DB_HOST" -P"$DB_PORT" --silent; do
  echo "Waiting for database connection..."
  sleep 2
done

echo "MySQL is ready! Starting the app..."

exec "$@"
