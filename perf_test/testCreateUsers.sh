#!/bin/bash
set -e

echo "Создание пользователей"
echo "POST /api/auth/signup"

TARGET_URL="http://localhost:5458"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RESULTS_DIR="results"
mkdir -p "$RESULTS_DIR"

echo ""
echo "Тест 1: Измерение RPS"
echo ""

wrk -t4 -c100 -d30s \
    -s ./create_users.lua \
    --timeout 10s \
    --latency \
    "$TARGET_URL" 2>&1 | tee "$RESULTS_DIR/create_test_$TIMESTAMP.txt"


echo ""
echo "Тест 2: Создание 100,000 пользователей"
echo ""

USER_COUNT=100000
BATCH_SIZE=1000

for ((i=1; i<=$USER_COUNT; i+=$BATCH_SIZE)); do
    END=$((i+BATCH_SIZE-1))
    if [ $END -gt $USER_COUNT ]; then
        END=$USER_COUNT
    fi
    
    echo "Создаю пользователей $i-$END из $USER_COUNT..."
    
    timeout 60 wrk -t2 -c50 -d60s \
        -s ./create_users.lua \
        --timeout 10s \
        "$TARGET_URL" > /dev/null 2>&1
    
    sleep 2
done

echo ""
echo "Создано 100,000 пользователей"
echo "Результаты сохранены в: $RESULTS_DIR/create_test_$TIMESTAMP.txt"


echo ""
echo "Статистика из wrk (30s тест):"
grep -A5 "Thread Stats" "$RESULTS_DIR/create_test_$TIMESTAMP.txt"
grep -A2 "Latency Distribution" "$RESULTS_DIR/create_test_$TIMESTAMP.txt" || true
grep "Requests/sec" "$RESULTS_DIR/create_test_$TIMESTAMP.txt"
grep "Transfer/sec" "$RESULTS_DIR/create_test_$TIMESTAMP.txt"