set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"
ENV_FILE="$PROJECT_ROOT/.env"

if [ -f "$ENV_FILE" ]; then
    source "$ENV_FILE"
else
    echo "ERROR: .env file not found"
    exit 1
fi

: ${DB_NAME:?"DB_NAME must be set"}
: ${APP_DB_PASSWORD:?"APP_DB_PASSWORD must be set"}


psql -v ON_ERROR_STOP=0 -d "$DB_NAME" -c "CREATE USER ddfilms_user;"

psql -v ON_ERROR_STOP=1 -d "$DB_NAME" -c "ALTER USER ddfilms_user WITH PASSWORD '$APP_DB_PASSWORD';"

psql -v ON_ERROR_STOP=1 -d "$DB_NAME" -f ../grantUserDB.sql
