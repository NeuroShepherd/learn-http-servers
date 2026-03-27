

source .env
cd sql/schema
goose postgres $DB_URL $1
