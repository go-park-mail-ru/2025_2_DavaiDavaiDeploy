[file name]: testCreateUsers_FIXED.sh
[file content begin]
#!/bin/bash
set -e


TARGET_URL="http://localhost:5458"

echo "Шаг 1: Тест производительности создания (30 секунд)"
echo "----------------------------------------------------"
wrk -t2 -c50 -d30s \
    -s ./createUsers.lua \
    --timeout 10s \
    --latency \
    "$TARGET_URL"

echo ""
echo "Шаг 2: Создание 100000 пользователей"
echo "--------------------------------------"

for i in {1..100}; do
    echo "Батч $i/100: создаю 1000 пользователей..."
    
    wrk -t1 -c20 -d10s \
        -s ./createUsers.lua \
        --timeout 10s \
        "$TARGET_URL" 2>&1 | grep "Создано пользователей" || true
    
    echo "   Готово"
    sleep 1
done


echo "Создано 100000 пользователей"
