set -e

CONF_DIR="/etc/postgresql/16/main/conf.d"
MAIN_CONF="/etc/postgresql/16/main/postgresql.conf"

sudo mkdir -p "$CONF_DIR"

sudo cp ../config/postgresql.conf "$CONF_DIR/kinopoisk.conf"

if ! grep -q "kinopoisk.conf" "$MAIN_CONF"; then
    echo "" | sudo tee -a "$MAIN_CONF"
    echo "# Kinopoisk application configuration" | sudo tee -a "$MAIN_CONF"
    echo "include = 'conf.d/kinopoisk.conf'" | sudo tee -a "$MAIN_CONF"
fi

sudo systemctl restart postgresql
