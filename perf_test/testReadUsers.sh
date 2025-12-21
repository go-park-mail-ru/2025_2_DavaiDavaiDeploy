#!/bin/bash
set -e

echo "Endpoint: GET /api/users/{id}"
echo ""

source ../.env

TARGET_URL="http://localhost:5458"

echo "Шаг 1: Получаем ID пользователей из базы"
./getUsersID.sh

echo ""
echo "Шаг 2: Нагрузочный тест чтения"
echo "Конфигурация: wrk -t2 -c50 -d30s"
echo ""

wrk -t2 -c50 -d30s \
    -s ./readUsers.lua \
    --timeout 10s \
    --latency \
    "$TARGET_URL"
