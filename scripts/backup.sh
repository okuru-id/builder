#!/bin/sh
# PostgreSQL backup via pg_dump into a timestamped, gzip-compressed file.
# Retains the last 7 days. Runs inside the app container.
set -e

BACKUP_DIR="${BACKUP_DIR:-/var/backups/okuru}"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p "$BACKUP_DIR"

OUT="$BACKUP_DIR/okuru_$DATE.sql"
PGPASSWORD="${DB_PASSWORD:-password}" pg_dump \
	-h "${DB_HOST:-postgres}" \
	-p "${DB_PORT:-5432}" \
	-U "${DB_USERNAME:-postgres}" \
	-d "${DB_DATABASE:-okuruid}" \
	--no-owner \
	--clean \
	>"$OUT"
gzip -f "$OUT"
find "$BACKUP_DIR" -name "okuru_*.sql.gz" -mtime +7 -delete
echo "Backup complete: $(basename "$OUT").gz"
