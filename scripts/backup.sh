#!/bin/sh
# SQLite backup with VACUUM into a timestamped, gzip-compressed file.
# Retains the last 7 days. Runs inside the app container.
set -e

DB_PATH="${DB_DATABASE:-okuru.db}"
# Resolve relative paths against the app working directory.
case "$DB_PATH" in
	/*) ;;
	*) DB_PATH="/opt/okuru/$DB_PATH" ;;
esac

BACKUP_DIR="${BACKUP_DIR:-/var/backups/okuru}"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p "$BACKUP_DIR"

if [ ! -f "$DB_PATH" ]; then
	echo "backup: database not found at $DB_PATH" >&2
	exit 1
fi

OUT="$BACKUP_DIR/okuru_$DATE.db"
sqlite3 "$DB_PATH" ".backup '$OUT'"
gzip -f "$OUT"
find "$BACKUP_DIR" -name "okuru_*.db.gz" -mtime +7 -delete
echo "Backup complete: $(basename "$OUT").gz"
