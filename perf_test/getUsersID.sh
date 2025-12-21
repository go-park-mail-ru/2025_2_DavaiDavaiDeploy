#!/bin/bash
set -e

echo "=== Получаем 100k ID пользователей из базы ==="

source ../.env

export PGPASSWORD=$DB_PASS

psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t \
  -c "SELECT id FROM user_table ORDER BY created_at DESC LIMIT 100000;" \
  | grep -v '^[[:space:]]*$' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//' \
  > user_ids.txt

COUNT=$(wc -l < user_ids.txt)
echo "Получено $COUNT ID пользователей"
echo "Файл: user_ids.txt"