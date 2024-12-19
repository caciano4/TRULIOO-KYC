#!/bin/bash

# Function to build the migrate command
build_command() {
  CMD="-path /app/migrations -database 'postgres://${DB_USER}:${DB_PASSWORD}@db_kyc:5432/${DB_NAME}?sslmode=disable'"
}

# Prompt user for migration options
echo "We have some options to migrate, Select from 1 to 4:"
echo "1: Migrate all"
echo "2: Rollback a specific number of migrations"
echo "3: Rollback all migrations"
echo "4: Check the current migration version"

# Read user input
read -p "Enter your choice: " number

# Build the base command
build_command

echo $MIGRATION_ACTION

# Handle user choice
case $MIGRATION_ACTION in
  1)
    echo "Running: migrate $CMD up"
    migrate -path /app/migrations -database 'postgres://caciano4:123caciano@db_kyc:5432/trullio?sslmode=disable' up
    ;;
  2)
    echo "Running: migrate $CMD down"
    migrate -path /app/migrations -database 'postgres://caciano4:123caciano@db_kyc:5432/trullio?sslmode=disable' down -all
    ;;
  3)
    echo "Running: migrate $CMD down -all"
    migrate -path /app/migrations -database 'postgres://caciano4:123caciano@db_kyc:5432/trullio?sslmode=disable' down $MIGRATION_STEPS
    ;;
  4)
    echo "Running: $CMD version"
    migrate -path /app/migrations -database 'postgres://caciano4:123caciano@db_kyc:5432/trullio?sslmode=disable' version
    ;;
  *)
    echo "Invalid optinos"
esac

# Keep the container running
exec tail -f /dev/null