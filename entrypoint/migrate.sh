#!/bin/bash

echo "Migration action: $MIGRATION_ACTION"

# Wait for database to be ready
echo "Waiting for database to be ready..."
until docker exec db_kyc pg_isready -U ${DB_USER} -d ${DB_NAME} > /dev/null 2>&1; do
  sleep 1
done
echo "Database is ready!"

# Function to check if table exists
table_exists() {
  local table_name=$1
  docker exec db_kyc psql -U ${DB_USER} -d ${DB_NAME} -tAc "SELECT 1 FROM information_schema.tables WHERE table_name='$table_name';" | grep -q 1
}

# Handle migration action
case $MIGRATION_ACTION in
  up)
    echo "Running migrations up..."
    for file in /app/migrations/*.up.sql; do
      if [ -f "$file" ]; then
        filename=$(basename "$file")
        table_name=$(echo "$filename" | grep -o 'create_[^_]*' | sed 's/create_//')
        
        if [ -n "$table_name" ] && table_exists "$table_name"; then
          echo "Skipping: $filename (table '$table_name' already exists)"
        else
          echo "Executing: $filename"
          docker exec -i db_kyc psql -U ${DB_USER} -d ${DB_NAME} < "$file"
        fi
      fi
    done
    echo "All migrations completed successfully!"
    ;;
  down)
    echo "Running migrations down..."
    STEPS=${MIGRATION_STEPS:-1}
    files=($(ls -r /app/migrations/*.down.sql 2>/dev/null))
    for ((i=0; i<$STEPS && i<${#files[@]}; i++)); do
      file=${files[$i]}
      if [ -f "$file" ]; then
        echo "Executing: $(basename $file)"
        docker exec -i db_kyc psql -U ${DB_USER} -d ${DB_NAME} < "$file"
      fi
    done
    echo "Migration rollback completed!"
    ;;
  version)
    echo "Checking database tables..."
    docker exec db_kyc psql -U ${DB_USER} -d ${DB_NAME} -c "\dt"
    ;;
  *)
    echo "Invalid migration action: $MIGRATION_ACTION"
    echo "Valid actions: up, down, version"
    exit 1
esac