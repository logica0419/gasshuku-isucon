#!/bin/bash -e
cd "$(dirname "$0")"

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-3306}
DB_USER=${DB_USER:-isucon}
DB_PASS=${DB_PASS:-isucon}
DB_NAME=${DB_NAME:-isulibrary}

cat 0_schema.sql 1_data.sql | mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME"
