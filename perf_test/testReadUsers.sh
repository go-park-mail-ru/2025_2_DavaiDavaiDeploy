#!/bin/bash
set -e

echo "=========================================="
echo "ТЕСТ ЧТЕНИЯ ПОЛЬЗОВАТЕЛЕЙ"
echo "GET /api/users/{id}"


source ../.env

TARGET_URL="http://localhost:5458"
RESULTS_DIR="results"
mkdir -p "$RESULTS_DIR"

echo ""
echo "Шаг 1: Получаем ID пользователей"
./get_user_ids.sh


echo ""
echo "Шаг 2: Нагрузочный тест чтения"
echo ""

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
wrk -t4 -c100 -d30s \
    -s ./readUsers.lua \
    --timeout 10s \
    --latency \
    "$TARGET_URL" 2>&1 | tee "$RESULTS_DIR/read_test_$TIMESTAMP.txt"


echo ""
echo "РЕЗУЛЬТАТЫ"
echo "Дата: $(date)"
echo "Endpoint: GET /api/users/{id}"
echo "Всего ID в файле: $(wc -l < user_ids.txt)"
echo ""
echo "Статистика:"
echo "----------------------------------------"
grep -A5 "Thread Stats" "$RESULTS_DIR/read_test_$TIMESTAMP.txt"
echo ""
grep "Requests/sec" "$RESULTS_DIR/read_test_$TIMESTAMP.txt"
grep "Transfer/sec" "$RESULTS_DIR/read_test_$TIMESTAMP.txt"
echo "----------------------------------------"